package algos

import (
    "github.com/roboticslab-uc3m/goraspio/ops"
)

type Pid struct {
    kp, ki, kd float64
    dt float64
    alpha float64
    prevValue float64
    integral float64
    derivative float64
    integralBounds [2]float64
}

func NewPid(kp, ki, kd, dt, alpha float64, integralBounds [2]float64) *Pid {
    return &Pid{kp, ki, kd, dt, alpha, 0, 0, 0, integralBounds}
}

func (p *Pid) Compute(value float64) float64 {
    p.integral += value * p.dt
    p.integral = ops.Clip(p.integral, p.integralBounds[0], p.integralBounds[1])

    unfiltDerivative := (value - p.prevValue) / p.dt
    p.derivative = p.alpha * unfiltDerivative + (1 - p.alpha) * p.derivative

    output := p.kp * value + p.ki * p.integral + p.kd * p.derivative

    p.prevValue = value

    return output
}
