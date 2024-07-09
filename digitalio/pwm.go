package digitalio

import (
    "fmt"
    "github.com/stianeikeland/go-rpio/v4"
)

type Pwm struct {
    pin rpio.Pin
    cycleLen uint32
}

func NewPwm(pinNo int, freq int) (Pwm, error) {
    clock := 64_000
    cycleLen := uint32(clock / freq)

    pin := rpio.Pin(pinNo)

    pin.Mode(rpio.Pwm)
    pin.Freq(clock)

    pin.Write(rpio.Low)
    pin.DutyCycle(0, cycleLen)

    return Pwm{pin, cycleLen}, nil
}

func (p Pwm) Write(dutyCycle int) error {
    if dutyCycle < 0 || dutyCycle > 100 {
        return fmt.Errorf("Duty cycle must be between 0 and 100")
    }

    // TMP
    if dutyCycle == 100 {
        dutyCycle = 99
    }

    dutyFreq := uint32(float32(dutyCycle) / 100.0 * float32(p.cycleLen))
    p.pin.DutyCycle(dutyFreq, p.cycleLen)

    return nil
}

func (p Pwm) Close() {
    p.pin.Write(rpio.Low)
    p.pin.DutyCycle(0, p.cycleLen)
}
