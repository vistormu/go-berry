package gpio

import (
    "time"

    "github.com/vistormu/goraspio/errors"
    "github.com/vistormu/goraspio/warnings"
    "github.com/vistormu/goraspio/num"
)

type PwmMode uint8
const (
    MarkSpace PwmMode = iota
    Balanced
)

const (
    pwmSamplingRate = 64_000 // for now maintain it as a constant
)

type Pwm struct {
    pin uint8
    mode PwmMode
    cycleLen uint32
    frequency int
}

func NewPwm(pinNo int) (*Pwm, error) {
    pin := uint8(pinNo)

    switch pinNo {
    case 12, 13, 40, 41, 45:
        setPinMode(pin, 4)
    case 18, 19:
        setPinMode(pin, 2)
    default:
        return nil, errors.New(errors.PWM_PIN, pinNo)
    }

    // freq
    initialFreq := 500
    cycleLen := uint32(pwmSamplingRate / initialFreq)

    pwm := &Pwm{
        pin: pin,
        cycleLen: cycleLen,
        mode: MarkSpace,
        frequency: initialFreq,
    }
    
    pwm.setSamplingRate(pwmSamplingRate)
    pwm.write(0)

    return pwm, nil
}

func (p *Pwm) SetFrequency(freq int) {
    if freq <= 0 || freq > pwmSamplingRate{
        prevFreq := freq
        freq = num.Clip(freq, 0, pwmSamplingRate)
        warnings.New(warnings.FREQUENCY, prevFreq, freq)
    }

    p.cycleLen = uint32(pwmSamplingRate / freq)
    p.frequency = freq
}

func (p *Pwm) Frequency() int {
    return p.frequency
}

func (p *Pwm) write(dutyCycle int) {
    dutyLen := uint32(float64(dutyCycle) / 100.0 * float64(p.cycleLen))

	const pwmCtlReg = 0
	var (
		pwmDatReg uint
		pwmRngReg uint
		shift     uint // offset inside ctlReg
	)

	switch p.pin {
	case 12, 18, 40: // channel pwm0
		pwmRngReg = 4
		pwmDatReg = 5
		shift = 0
	case 13, 19, 41, 45: // channel pwm1
		pwmRngReg = 8
		pwmDatReg = 9
		shift = 8
	}

	const ctlMask = 255 // ctl setting has 8 bits for each channel
	const pwen = 1 << 0 // enable pwm
	var msen uint32 = 0
	if p.mode == MarkSpace {
		msen = 1 << 7
	}

	pwmMem[pwmCtlReg] = pwmMem[pwmCtlReg]&^(ctlMask<<shift) | msen<<shift | pwen<<shift

	// set duty cycle
	pwmMem[pwmDatReg] = dutyLen
	pwmMem[pwmRngReg] = p.cycleLen

	time.Sleep(time.Microsecond * 10)
}
func (p *Pwm) Write(dutyCycle int) {
    if dutyCycle < 0 || dutyCycle > 100 {
        prevDutyCycle := dutyCycle
        dutyCycle = num.Clip(dutyCycle, 0, 100)
        warnings.New(warnings.DUTY_CYCLE, prevDutyCycle, dutyCycle)
    }

    p.write(dutyCycle)
}

func (p *Pwm) Close() {
    p.write(0)  
}

func (p *Pwm) start() {
	const pwmCtlReg = 0
	const pwen = 1
	pwmMem[pwmCtlReg] |= pwen<<8 | pwen
}

func (p *Pwm) stop() {
	const pwmCtlReg = 0
	const pwen = 1
	pwmMem[pwmCtlReg] &^= pwen<<8 | pwen
}

func (p *Pwm) setSamplingRate(freq int) {
	sourceFreq := 19200000 // oscilator frequency
	if isBCM2711() {
		sourceFreq = 52000000
	}
	const divMask = 4095 // divi and divf have 12 bits each

	divi := uint32(sourceFreq / freq)
	divf := uint32(((sourceFreq % freq) << 12) / freq)

	divi &= divMask
	divf &= divMask

	clkCtlReg := 28
	clkDivReg := 28
    clkCtlReg += 12
    clkDivReg += 13

    p.stop()
    defer p.start()

	mash := uint32(1 << 9) // 1-stage MASH
	if divi < 2 || divf == 0 {
		mash = 0
	}

	memlock.Lock()
	defer memlock.Unlock()

	const PASSWORD = 0x5A000000
	const busy = 1 << 7
	const enab = 1 << 4
	const src = 1 << 0 // oscilator

	clkMem[clkCtlReg] = PASSWORD | (clkMem[clkCtlReg] &^ enab) // stop gpio clock (without changing src or mash)
	for clkMem[clkCtlReg]&busy != 0 {
		time.Sleep(time.Microsecond * 10)
	} // ... and wait for not busy

	clkMem[clkCtlReg] = PASSWORD | mash | src          // set mash and source (without enabling clock)
	clkMem[clkDivReg] = PASSWORD | (divi << 12) | divf // set dividers

	// mash and src can not be changed in same step as enab, to prevent lock-up and glitches
	time.Sleep(time.Microsecond * 10) // ... so wait for them to take effect

	clkMem[clkCtlReg] = PASSWORD | mash | src | enab // finally start clock
}
