package sensor

import (
    "fmt"
	"github.com/vistormu/goraspio/digitalio"
)

type Ems20 struct {
    spi digitalio.Spi
}

func NewEms20(chipSelectPinNo int) (Ems20, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo)
    if err != nil {
        return Ems20{}, fmt.Errorf("error opening communication channel\n%v", err)
    }

    lc :=  Ems20{
        spi: spi,
    }

    return lc, nil
}

func (lc Ems20) read() (int, error) {
    // read bytes
    data, err := lc.spi.Read()
    if err != nil {
        return -1, fmt.Errorf("error reading channel\n%v", err)
    }

    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)

    return value, nil
}

func (lc Ems20) Read() (float64, error) {
    value, err := lc.read()
    if err != nil {
        return -1.0, fmt.Errorf("error reading value\n%v", err)
    }

    load := (float64(value) / 4095) * 50

    return load, nil
}

func (s Ems20) Close() error {
    s.spi.Close()

    return nil
}
