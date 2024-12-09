// polarity 0
// phase 1
// mode 1
// spi
// 14-bits
package sensor

import (
	"github.com/vistormu/goraspio/gpio"
)

type As5048a struct {
    spi *gpio.Spi 
    offset int
    prevData int
    resetCount int
}

func NewAs5048a(chipSelectNo int) (*As5048a, error) {
    spi, err := gpio.NewSpi(chipSelectNo, 0, 1, 10_000)
    if err != nil {
        return nil, err
    }

    s := &As5048a{spi, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, err
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5048a) read() (int, error) {
    data := s.spi.Read(2)

    value := (uint16(data[0]) << 8) | uint16(data[1])
    angleValue := value & 0x3FFF

    return int(angleValue), nil
}

func (s *As5048a) Read() (float64, error) {
    data, err := s.read()
    if err != nil {
        return -1, err
    }

    degreeValue := float64(data) / 16384 * 360
    
    return degreeValue, nil
}

func (s *As5048a) Close() error {
    s.spi.Close()

    return nil
}

