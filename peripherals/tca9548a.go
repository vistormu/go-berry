package peripherals

import (
	"github.com/vistormu/go-berry/comms"
)

type Tca9548a struct {
	i2c *comms.I2C
}

func NewTca9548a(address uint8, bus int) (*Tca9548a, error) {
	i2c, err := comms.NewI2C(address, bus)
	if err != nil {
		return nil, err
	}

	return &Tca9548a{
		i2c: i2c,
	}, nil
}

func (s *Tca9548a) Select(channel uint8) error {
	err := s.i2c.Write(0, 1<<channel)
	if err != nil {
		return err
	}

	return nil
}

func (s *Tca9548a) Close() error {
	return s.i2c.Close()
}
