package sensor

import (
    "fmt"
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
        return nil, fmt.Errorf("error opening communication channel\n%v", err)
    }

    s := &As5048b{i2cChannel, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, fmt.Errorf("error reading initial value\n%v", err)
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5048b) read() (int, error) {
    highByte, err := s.i2cChannel.ReadRegU8(0xFF)
    if err != nil {
        return -1, fmt.Errorf("error reading 0xFF channel\n%v", err)
    }
    lowByte, err := s.i2cChannel.ReadRegU8(0x0FE)
    if err != nil {
        return -1, fmt.Errorf("error reading 0x0FE channel\n%v", err)
    }

    value := (uint16(highByte) << 6) | (uint16(lowByte) & 0x3F)

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

