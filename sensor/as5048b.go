package sensor

import (
    "github.com/d2r2/go-i2c"
    "github.com/d2r2/go-logger"
)

type As5048b struct {
    i2cChannel *i2c.I2C
    offset int
    prevData int
    resetCount int
}

func NewAs5048b(address byte, line int) (*As5048b, error) {
    logger.ChangePackageLogLevel("i2c", logger.FatalLevel)

    i2cChannel, err := i2c.NewI2C(address, line)
    if err != nil {
        return nil, err
    }

    s := &As5048b{i2cChannel, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, err
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5048b) read() (int, error) {
    highByte, err := s.i2cChannel.ReadRegU8(0x00)
    if err != nil {
        return -1, err
    }
    lowByte, err := s.i2cChannel.ReadRegU8(0x01)
    if err != nil {
        return -1, err
    }

    value := (uint16(highByte) << 8) | uint16(lowByte)
    value = value & 0x3FFF

    return int(value), nil
}

func (s *As5048b) Read() (float64, error) {
    data, err := s.read()
    if err != nil {
        return -1, err
    }

    degreeValue := float64(data) / 16384 * 360
    
    return degreeValue, nil
}

func (s *As5048b) Close() error {
    err := s.i2cChannel.Close()
    if err != nil {
        return err
    }

    return nil
}

