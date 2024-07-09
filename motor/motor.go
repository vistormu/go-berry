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

func (m Motor) Write(posError float64) (int, error) {
    // direction and PWM
    pwmValue := 0
    if posError < -TOLERANCE {
        m.direction.Write(digitalio.High)
        pwmValue = 100
    } else if posError > TOLERANCE {
        m.direction.Write(digitalio.Low)
        pwmValue = 100
    } else {
        pwmValue = 0
    }

    err := m.pwm.Write(pwmValue)
    if err != nil {
        return 0, err
    }

    return pwmValue, nil
}

func (m Motor) Close() {
    m.pwm.Close()
    m.direction.Close()
}
