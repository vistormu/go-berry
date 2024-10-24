package actuator

import (
    "math"

    "github.com/roboticslab-uc3m/goraspio/digitalio"
    "github.com/roboticslab-uc3m/goraspio/controller"
    "github.com/roboticslab-uc3m/goraspio/utils"
)

const (
    TOLERANCE = 0.05 // mm
)

type Motor struct {
    pwm digitalio.Pwm
    direction digitalio.DigitalOut
    pid *controller.Pid
}


func New(pwmPinNo, freq, directionPinNo int) (Motor, error) {
    // pwm
    pwm, err := digitalio.NewPwm(pwmPinNo, freq)
    if err != nil {
        return Motor{}, err
    }
    
    // direction
    direction := digitalio.NewDigitalOut(directionPinNo, digitalio.Low)

    // pid
    pid := controller.NewPid(100, 0, 0, 0, [2]float64{-1, 1})

    return Motor{pwm, direction, pid}, nil
}

func (m Motor) Write(posError float64, dt float64) (int, error) {
    // pwm value
    pwmValue := m.pid.Compute(float64(posError), float64(dt))
    pwmValue = utils.Clip(pwmValue, -100, 100)

    sign := math.Signbit(float64(pwmValue))
    value := float64(math.Abs(float64(pwmValue)))

    // direction
    if sign {
        m.direction.Write(digitalio.High) // negative error
    } else {
        m.direction.Write(digitalio.Low) // positive error
    }

    // write
    err := m.pwm.Write(int(value))
    if err != nil {
        return 0, err
    }

    return int(pwmValue), nil
}

func (m Motor) WriteRaw(pwmValue int, direction digitalio.PinState) error {
    m.direction.Write(direction)
    err := m.pwm.Write(pwmValue)
    if err != nil {
        return err
    }

    return nil
}

func (m Motor) Close() {
    m.pwm.Close()
    m.direction.Close()
}
