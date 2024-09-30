package sensor

import (
    "math"
    "github.com/d2r2/go-i2c"
    "github.com/d2r2/go-logger"
)

const (
    RESET_THRESH = 0.4
    MAX_VALUE = 4095
    STEP_TO_MM = 0.000488
)

type As5311 struct {
    i2cChannel *i2c.I2C
    offset int
    prevData int
    resetCount int
    prevValue float64
}

func NewAs5311(address byte, line int) (*As5311, error) {
    logger.ChangePackageLogLevel("i2c", logger.FatalLevel)

    i2cChannel, err := i2c.NewI2C(address, line)
    if err != nil {
        return nil, err
    }

    s := &As5311{i2cChannel, 0, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, err
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5311) read() (int, error) {
    highByte, err := s.i2cChannel.ReadRegU8(0x00)
    if err != nil {
        return -1, err
    }
    lowByte, err := s.i2cChannel.ReadRegU8(0x01)
    if err != nil {
        return -1, err
    }

    value := (int(highByte) << 4) | (int(lowByte) >> 4)

    return value, nil
}

func (s *As5311) Read() (float64, error) {
    data, err := s.read()
    if err != nil {
        return s.prevValue, err
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
    err := s.i2cChannel.Close()
    if err != nil {
        return err
    }

    return nil
}

