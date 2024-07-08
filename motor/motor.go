package motor

import (
	"goraspio/digitalio"
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

func (m Motor) Write(reference, position float64) error {
    // error
    posError := reference - position

    // direction and PWM
    var err error
    if posError < TOLERANCE {
        m.direction.Write(digitalio.Low)
        err = m.pwm.Write(100)
    } else if posError > TOLERANCE {
        m.direction.Write(digitalio.High)
        err = m.pwm.Write(100)
    } else {
        err = m.pwm.Write(0)
    }

    if err != nil {
        return err
    }

    return nil
}

func (m Motor) Close() {
    m.pwm.Close()
    m.direction.Close()
}
