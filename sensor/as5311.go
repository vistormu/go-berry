// As5311
// 488 nm resolution
// 12-bits
// SSI
package sensor

import (
    "math"
    "fmt"
	"github.com/vistormu/goraspio/gpio"
)

// const (
//     RESET_THRESH = 0.4
//     MAX_VALUE = 4095
//     STEP_TO_MM = 0.000488
// )

type As5311 struct {
    spi *gpio.Spi 
    offset int
    prevData int
    resetCount int
    prevValue float64
}

func NewAs5311(chipSelectNo int) (*As5311, error) {
    spi, err := gpio.NewSpi(chipSelectNo, 1, 1, 10_000)
    if err != nil {
        return nil, fmt.Errorf("error opening communication channel\n%v", err)
    }

    s := &As5311{spi, 0, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, fmt.Errorf("error reading initial value\n%v", err)
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5311) read() (int, error) {
    data := s.spi.Read(2)

    value := (int(data[0]) << 4) | int(data[1] >> 4)

    return int(value), nil
}

func (s *As5311) Read() (float64, error) {
    data, err := s.read()
    if err != nil {
        return s.prevValue, fmt.Errorf("error reading value\n%v", err)
    }

    diff := float64(data - s.prevData)
    change := float64(MAX_VALUE)*(1-RESET_THRESH)
    if diff < 0 && math.Abs(diff) > change {
        s.resetCount++
    }
    if diff > 0 && math.Abs(diff) > change {
        s.resetCount--
    }

    s.prevData = data
    output := data - s.offset + s.resetCount * (MAX_VALUE + 1)

    position := -float64(output)*STEP_TO_MM
    s.prevValue = position

    return position, nil
}

func (s *As5311) Close() error {
    s.spi.Close()

    return nil
}
