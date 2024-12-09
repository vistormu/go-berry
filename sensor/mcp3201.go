package sensor

import (
	"github.com/vistormu/goraspio/gpio"
)

type Mcp3201 struct {
    spi *gpio.Spi
    vRef float64
}

func NewMcp3201(vRef float64, chipSelectPinNo int) (Mcp3201, error) {
    spi, err := gpio.NewSpi(chipSelectPinNo, 0, 0, 16_000) 
    if err != nil {
        return Mcp3201{}, err
    }
    
    return Mcp3201{
        spi: spi,
        vRef: vRef,
    }, nil
}

func (m Mcp3201) Read() float64 {
    data := m.spi.Read(2)
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float64(value) / 4095) * m.vRef

    return voltage
}

func (m Mcp3201) Close() {
    m.spi.Close()
}
