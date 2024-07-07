package controller

type Controller interface {
    Compute(err, dt float32) float32
}
