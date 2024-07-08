package motor

import (
	"errors"
	"goraspio/digitalio"
	"math"
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
    if value < -1.0 && value > 1.0 {
        return errors.New("Motor value must be between -1.0 and 1.0") 
    }

    // direction
    if math.Signbit(value) { // negative
        m.direction.Write(digitalio.High)
    } else {
        m.direction.Write(digitalio.Low)
    }

    // pwm
    pwmValue := int(math.Abs(value*100))
    m.pwm.Write(pwmValue)

    return nil
}

func (m Motor) Close() {
    m.pwm.Close()
    m.direction.Close()
}
