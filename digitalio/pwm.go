package digitalio

import (
    "fmt"
    "github.com/stianeikeland/go-rpio/v4"
)

const (
    samplingRate = 64_000
)

type Pwm struct {
    pin rpio.Pin
    cycleLen int
    frequency int
}

func NewPwm(pinNo int) *Pwm {
    // freq
    initialFreq := 500
    cycleLen := int(samplingRate / initialFreq)

    // pin
    pin := rpio.Pin(pinNo)
    pin.Mode(rpio.Pwm)
    pin.Freq(samplingRate)

    pin.Write(rpio.Low)
    pin.DutyCycle(0, uint32(cycleLen))

    return &Pwm{
        pin: pin,
        cycleLen: cycleLen,
    }
}

func (p *Pwm) SetFrequency(freq int) error {
    if freq <= 0 {
        return fmt.Errorf("frequency must be positive, got: %d", freq)
    }

    p.cycleLen = int(samplingRate / freq)
    p.frequency = freq

    return nil
}

func (p *Pwm) Frequency() int {
    return p.frequency
}

func (p *Pwm) Write(dutyCycle int) error {
    if dutyCycle < 0 || dutyCycle > 100 {
        return fmt.Errorf("duty cycle must be between 0 and 100, got: %d", dutyCycle)
    }

    duty := uint32(float64(dutyCycle) / 100.0 * float64(p.cycleLen))
    p.pin.DutyCycle(duty, uint32(p.cycleLen))

    return nil
}

func (p *Pwm) Close() {
    p.pin.Write(rpio.Low)
    p.pin.DutyCycle(0, uint32(p.cycleLen))
}
