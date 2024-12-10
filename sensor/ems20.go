package sensor

import (
	"github.com/vistormu/goraspio/gpio"
)

type Ems20 struct {
    spi *gpio.Spi
}

func NewEms20(chipSelectPinNo int) (*Ems20, error) {
    spi, err := gpio.NewSpi(chipSelectPinNo, 0, 0, 10_000)
    if err != nil {
        return nil, err
    }

    return &Ems20{
        spi: spi,
    }, nil
}

func (lc Ems20) read() (int, error) {
    data, err := lc.spi.Read(2)
    if err != nil {
        return -1, err
    }

    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)

    return value, nil
}

func (lc Ems20) Read() (float64, error) {
    value, err := lc.read()
    if err != nil {
        return -1.0, err
    }

    load := (float64(value) / 4095) * 50

    return load, nil
}

func (s Ems20) Close() error {
    err := s.spi.Close()
    if err != nil {
        return err
    }
    
    return nil
}
