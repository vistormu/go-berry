package sensor

type Sensor interface {
    Read() (any, error)
    Close() error
}
