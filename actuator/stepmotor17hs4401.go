package actuator

import (
    "github.com/vistormu/goraspio/digitalio"
    "github.com/vistormu/goraspio/ops"
)

type StepMotor17hs4401 struct {
    pwm *digitalio.Pwm
    direction digitalio.DigitalOut
    minFreq int
    maxFreq int
}

func NewStepMotor17hs4401(stepPinNo, directionPinNo, minFreq, maxFreq int) (*StepMotor17hs4401, error) {
    motor := &StepMotor17hs4401{
        pwm: digitalio.NewPwm(stepPinNo),
        direction: digitalio.NewDigitalOut(directionPinNo, digitalio.Low),
        minFreq: minFreq,
        maxFreq: maxFreq,
    }
    return motor, nil
}

func (m *StepMotor17hs4401) Write(value float64) error {
    speed := ops.Clip(int(value), -100, 100)
    frequency := ops.MapInterval(ops.Abs(speed), 0, 100, m.minFreq, m.maxFreq)

    if speed == 0 {
        return m.pwm.Write(0)
    }

    if speed > 0 {
        m.direction.Write(digitalio.Low)
    } else {
        m.direction.Write(digitalio.High)
    }

    err := m.pwm.SetFrequency(frequency)
    if err != nil {
        return err
    }
    
    return m.pwm.Write(50)
}

func (m *StepMotor17hs4401) Close() error {
    m.pwm.Close()
    m.direction.Close()
    
    return nil
}
