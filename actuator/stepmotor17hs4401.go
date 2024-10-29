package actuator

import (
    "math"
    "github.com/vistormu/goraspio/digitalio"
    "github.com/vistormu/goraspio/ops"
)

type StepMotor17hs4401 struct {
    pwm digitalio.Pwm
    direction digitalio.DigitalOut
}


func NewStepMotor17hs4401(pwmPinNo, freq, directionPinNo int) (StepMotor17hs4401, error) {
    pwm, err := digitalio.NewPwm(pwmPinNo, freq)
    if err != nil {
        return StepMotor17hs4401{}, err
    }
    
    direction := digitalio.NewDigitalOut(directionPinNo, digitalio.Low)

    return StepMotor17hs4401{pwm, direction}, nil
}

func (m StepMotor17hs4401) Write(value float64) error {
    if math.Signbit(value) {
        m.direction.Write(digitalio.High) // negative error
    } else {
        m.direction.Write(digitalio.Low) // positive error
    }

    // write
    pwmValue := int(math.Abs(value))
    pwmValue = ops.Clip(pwmValue, 0, 100)

    err := m.pwm.Write(pwmValue)
    if err != nil {
        return err
    }

    return nil
}

func (m StepMotor17hs4401) Close() error {
    m.pwm.Close()
    m.direction.Close()

    return nil
}
