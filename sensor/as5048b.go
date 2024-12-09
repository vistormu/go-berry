package sensor

import (
    "fmt"
    "github.com/vistormu/goraspio/gpio"
)

type As5048b struct {
    i2c *gpio.I2C
    offset int
    prevData int
    resetCount int
}

func NewAs5048b(address byte, line int) (*As5048b, error) {
    i2c, err := gpio.NewI2C(address, line)
    if err != nil {
        return nil, fmt.Errorf("error opening communication channel\n%v", err)
    }

    s := &As5048b{i2c, 0, 0, 0}

    s.offset, err = s.read()
    if err != nil {
        return nil, fmt.Errorf("error reading initial value\n%v", err)
    }
    s.prevData = s.offset

    return s, nil
}

func (s *As5048b) read() (int, error) {
    data, err := s.i2c.Read([]byte{0xFF, 0x0FE}, []int{1, 1})
    if err != nil {
        return -1, err
    }

    value := (uint16(data[0]) << 6) | (uint16(data[1]) & 0x3F)

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
    err := s.i2c.Close()
    if err != nil {
        return err
    }

    return nil
}

