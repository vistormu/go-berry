package sensor

import (
    "fmt"
	"github.com/roboticslab-uc3m/goraspio/digitalio"
)

type Mcp3201 struct {
    spi digitalio.Spi
    vRef float64
}

func NewMcp3201(vRef float64, chipSelectPinNo int) (Mcp3201, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo) 
    if err != nil {
        return Mcp3201{}, fmt.Errorf("error opening communication channel\n%v", err)
    }
    
    return Mcp3201{
        spi: spi,
        vRef: vRef,
    }, nil
}

func (m *Mcp3201) read() (int, error) {
    // read bytes
    data, err := m.spi.Read()
    if err != nil {
        return 0.0, fmt.Errorf("error reading channel\n%v", err)
    }
    
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)

    return value, nil
}

func (m *Mcp3201) Read() (float64, error) {
    value, err := m.read()
    if err != nil {
        return -1.0, fmt.Errorf("error reading value\n%v", err)
    }

    voltage := (float64(value) / 4095) * m.vRef

    return voltage, nil
}

func (m *Mcp3201) Close() error {
    m.spi.Close()

    return nil
}
