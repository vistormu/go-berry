package actuator

import (
    "github.com/vistormu/goraspio/gpio"
    "github.com/vistormu/goraspio/num"
)

type StepMotor17hs4401 struct {
    pwm *gpio.Pwm
    direction *gpio.DigitalOut
    minFreq int
    maxFreq int
}

func NewStepMotor17hs4401(stepPinNo, directionPinNo, minFreq, maxFreq int) (*StepMotor17hs4401, error) {
    pwm, err := gpio.NewPwm(stepPinNo)
    if err != nil {
        return nil, err
    }

    dir, err := gpio.NewDigitalOut(directionPinNo, gpio.Low)
    if err != nil {
        return nil, err
    }

    motor := &StepMotor17hs4401{
        pwm: pwm,
        direction: dir,
        minFreq: minFreq,
        maxFreq: maxFreq,
    }

    return motor, nil
}

func (m *StepMotor17hs4401) Write(value float64) error {
    speed := num.Clip(int(value), -100, 100)
    frequency := num.MapInterval(num.Abs(speed), 0, 100, m.minFreq, m.maxFreq)

    var err error
    if speed == 0 {
        err = m.pwm.Write(0)
        if err != nil {
            return err
        }
    }

    if speed > 0 {
        err = m.direction.Write(gpio.Low)
        if err != nil {
            return err
        }
    } else {
        err = m.direction.Write(gpio.High)
        if err != nil {
            return err
        }
    }

    m.pwm.SetFrequency(frequency)
    err = m.pwm.Write(50)
    if err != nil {
        return err
    }

    return nil
}

func (m *StepMotor17hs4401) Close() error {
    var err error
    err = m.pwm.Close()
    if err != nil {
        return err
    }

    err = m.direction.Close()
    if err != nil {
        return err
    }

    return nil
}
