package actuator

type Actuator interface {
    Write(value float64) error
    Close() error
}
