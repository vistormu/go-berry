// Nse5310
// 0.488 μm resolution
// 12-bits
// I2C
package peripherals

import (
    "math"
    "github.com/vistormu/go-berry/comms"
)

const (
    RESET_THRESH = 0.4
    MAX_VALUE = 4095
    STEP_TO_MM = 0.000488
)

type Nse5310 struct {
    i2cChannel *comms.I2C
    offset int
    prevData int
    resetCount int
    prevValue float64
}

func NewNse5310(address byte, line int) (*Nse5310, error) {
    i2cChannel, err := comms.NewI2C(address, line)
    if err != nil {
        return nil, err
    }

    s := &Nse5310{i2cChannel, 0, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, err
    }
    s.prevData = s.offset

    return s, nil
}

func (s *Nse5310) read() (int, error) {
    data, err := s.i2cChannel.Read([]byte{0x00, 0x01}, []int{1, 1})
    if err != nil {
        return -1, err
    }

    value := (int(data[0]) << 4) | (int(data[1]) >> 4)

    return value, nil
}

func (s *Nse5310) Read() (float64, error) {
    // read from i2c
    data, err := s.read()
    if err != nil {
        return s.prevValue, err
    }

    // calculate reset values
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

func (s *Nse5310) Close() error {
    err := s.i2cChannel.Close()
    if err != nil {
        return err
    }

    return nil
}
