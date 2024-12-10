package sensor

import (
	"github.com/vistormu/goraspio/gpio"
)

type Mcp3201 struct {
    spi *gpio.Spi
    vRef float64
}

func NewMcp3201(vRef float64, chipSelectPinNo int) (*Mcp3201, error) {
    spi, err := gpio.NewSpi(chipSelectPinNo, 0, 0, 1.6e6) 
    if err != nil {
        return nil, err
    }
    
    return &Mcp3201{
        spi: spi,
        vRef: vRef,
    }, nil
}

func (m Mcp3201) Read() (float64, error) {
    data, err := m.spi.Read(2)
    if err != nil {
        return 0, err
    }

    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float64(value) / 4095) * m.vRef

    return voltage, nil
}

func (m Mcp3201) Close() error {
    err := m.spi.Close()
    if err != nil {
        return err
    }

    return nil
}
