package peripherals

import (
	"github.com/vistormu/go-berry/comms"
)

type Mcp3204 struct {
    spi *comms.Spi
    vRef float64
}

func NewMcp3204(vRef float64, chipSelectPinNo int) (*Mcp3204, error) {
    spi, err := comms.NewSpi(chipSelectPinNo, 0, 0, 1.6e6) 
    if err != nil {
        return nil, err
    }
    
    return &Mcp3204{
        spi: spi,
        vRef: vRef,
    }, nil
}

func (m Mcp3204) Read(channel int) (float64, error) {
    cmd := make([]byte, 3)
    cmd[0] = 0x06 + (byte(channel) >> 2)
	cmd[1] = (byte(channel) & 0x03) << 6
	cmd[2] = 0x00

    data := make([]byte, 3)
    copy(data, cmd)
    m.spi.Exchange(data)

    value := ((int(data[1]) & 0x0F) << 8) | int(data[2])
    voltage := (float64(value) / 4095) * m.vRef

    return voltage, nil
}

func (m Mcp3204) Close() error {
    err := m.spi.Close()
    if err != nil {
        return err
    }

    return nil
}
