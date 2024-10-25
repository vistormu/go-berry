package actuator

import (
    "math"
    "github.com/roboticslab-uc3m/goraspio/digitalio"
)

const (
    TOLERANCE = 0.05 // mm
)

type Motor struct {
    pwm digitalio.Pwm
    direction digitalio.DigitalOut
}


func New(pwmPinNo, freq, directionPinNo int) (Motor, error) {
    pwm, err := digitalio.NewPwm(pwmPinNo, freq)
    if err != nil {
        return Motor{}, err
    }
    
    direction := digitalio.NewDigitalOut(directionPinNo, digitalio.Low)

    return Motor{pwm, direction}, nil
}

func (m Motor) Write(value float64) error {
    sign := math.Signbit(value)

    // direction
    if sign {
        m.direction.Write(digitalio.High) // negative error
    } else {
        m.direction.Write(digitalio.Low) // positive error
    }

    // write
    err := m.pwm.Write(int(value))
    if err != nil {
        return err
    }

    return nil
}

func (m Motor) Close() error {
    m.pwm.Close()
    m.direction.Close()

    return nil
}
