// Nse5310
// 0.488 μm resolution
// 12-bits
// I2C
package sensor

import (
    "math"
    "fmt"
    "github.com/d2r2/go-i2c"
    "github.com/d2r2/go-logger"
)

const (
    RESET_THRESH = 0.4
    MAX_VALUE = 4095
    STEP_TO_MM = 0.000488
)

type Nse5310 struct {
    i2cChannel *i2c.I2C
    offset int
    prevData int
    resetCount int
    prevValue float64
}

func NewNse5310(address byte, line int) (*Nse5310, error) {
    logger.ChangePackageLogLevel("i2c", logger.FatalLevel)

    i2cChannel, err := i2c.NewI2C(address, line)
    if err != nil {
        return nil, fmt.Errorf("error opening communication channel\n%v", err)
    }

    s := &Nse5310{i2cChannel, 0, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, fmt.Errorf("error reading initial value\n%v", err)
    }
    s.prevData = s.offset

    return s, nil
}

func (s *Nse5310) read() (int, error) {
    highByte, err := s.i2cChannel.ReadRegU8(0x00)
    if err != nil {
        return -1, fmt.Errorf("error reading 0x00 channel\n%v", err)
    }
    lowByte, err := s.i2cChannel.ReadRegU8(0x01)
    if err != nil {
        return -1, fmt.Errorf("error reading 0x01 channel\n%v", err)
    }

    value := (int(highByte) << 4) | (int(lowByte) >> 4)

    return value, nil
}

func (s *Nse5310) Read() (float64, error) {
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

func (s *Nse5310) Close() error {
    err := s.i2cChannel.Close()
    if err != nil {
        return err
    }

    return nil
}
