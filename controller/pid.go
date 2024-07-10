package controller

import (
    "goraspio/utils"
)

type Pid struct {
    kp, ki, kd float32
    alpha float32
    prevError float32
    integral float32
    derivative float32
    integralBounds [2]float32
}

func NewPid(kp, ki, kd, alpha float32, integralBounds [2]float32) *Pid {
    return &Pid{kp, ki, kd, alpha, 0, 0, 0, integralBounds}
}

func (p *Pid) Compute(err, dt float32) float32 {
    p.integral += err * dt
    p.integral = utils.Clip(p.integral, p.integralBounds[0], p.integralBounds[1])

    unfiltDerivative := (err - p.prevError) / dt
    p.derivative = p.alpha * unfiltDerivative + (1 - p.alpha) * p.derivative

    output := p.kp * err + p.ki * p.integral + p.kd * p.derivative

    p.prevError = err

    return output
}
