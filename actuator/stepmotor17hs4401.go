package actuator

import (
    "github.com/vistormu/goraspio/gpio"
    "github.com/vistormu/goraspio/num"
)

type StepMotor17hs4401 struct {
    pwm *gpio.Pwm
    direction gpio.DigitalOut
    minFreq int
    maxFreq int
}

func NewStepMotor17hs4401(stepPinNo, directionPinNo, minFreq, maxFreq int) (*StepMotor17hs4401, error) {
    pwm, err := gpio.NewPwm(stepPinNo)
    if err != nil {
        return nil, err
    }

    motor := &StepMotor17hs4401{
        pwm: pwm,
        direction: gpio.NewDigitalOut(directionPinNo, gpio.Low),
        minFreq: minFreq,
        maxFreq: maxFreq,
    }
    return motor, nil
}

func (m *StepMotor17hs4401) Write(value float64) {
    speed := num.Clip(int(value), -100, 100)
    frequency := num.MapInterval(num.Abs(speed), 0, 100, m.minFreq, m.maxFreq)

    if speed == 0 {
        m.pwm.Write(0)
    }

    if speed > 0 {
        m.direction.Write(gpio.Low)
    } else {
        m.direction.Write(gpio.High)
    }

    m.pwm.SetFrequency(frequency)
    m.pwm.Write(50)
}

func (m *StepMotor17hs4401) Close() {
    m.pwm.Close()
    m.direction.Close()
}
