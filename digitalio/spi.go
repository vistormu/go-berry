package digitalio

import (
    "fmt"
    "github.com/stianeikeland/go-rpio/v4"
)

type Spi struct {
    chipSelect DigitalOut
}

func NewSpi(chipSelectPinNo int) (Spi, error) {
    err := rpio.SpiBegin(rpio.Spi0)
    if err != nil {
        return Spi{}, err
    }
    rpio.SpiMode(0, 1)
    rpio.SpiSpeed(1_000_000)

    do := NewDigitalOut(chipSelectPinNo, High)

    return Spi{do}, nil
}

func (s Spi) Read() ([]byte, error) {
    s.chipSelect.Toggle()
    defer s.chipSelect.Toggle()

    data := rpio.SpiReceive(2)
    if len(data) != 2 {
        return nil, fmt.Errorf("Failed to read from SPI")
    }

    return data, nil
}

func (s Spi) Close() {
    rpio.SpiEnd(rpio.Spi0)
    s.chipSelect.Close()
}
