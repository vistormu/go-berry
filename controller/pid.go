package controller

import (
    "github.com/roboticslab-uc3m/goraspio/utils"
)

type Pid struct {
    kp, ki, kd float64
    alpha float64
    prevError float64
    integral float64
    derivative float64
    integralBounds [2]float64
}

func NewPid(kp, ki, kd, alpha float64, integralBounds [2]float64) *Pid {
    return &Pid{kp, ki, kd, alpha, 0, 0, 0, integralBounds}
}

func (p *Pid) Compute(err, dt float64) float64 {
    p.integral += err * dt
    p.integral = utils.Clip(p.integral, p.integralBounds[0], p.integralBounds[1])

    unfiltDerivative := (err - p.prevError) / dt
    p.derivative = p.alpha * unfiltDerivative + (1 - p.alpha) * p.derivative

    output := p.kp * err + p.ki * p.integral + p.kd * p.derivative

    p.prevError = err

    return output
}
