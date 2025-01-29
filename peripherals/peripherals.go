package peripherals

type Sensor interface {
    Read() (float64, error)
    Close() error
}
