package peripherals

import (
	"github.com/vistormu/go-berry/comms"
	"math"
)

type As5510 struct {
	i2cChannel *comms.I2C
	offset     int
	prevData   int
	resetCount int
	prevValue  float64

	resetThresh float64
	maxValue    int
	stepToMm    float64
}

func NewAs5510(address byte, line int) (*As5510, error) {
	i2cChannel, err := comms.NewI2C(address, line)
	if err != nil {
		return nil, err
	}

	s := &As5510{
		i2cChannel:  i2cChannel,
		resetThresh: 0.4,
		maxValue:    1024 - 1,
		stepToMm:    1.0,
	}

	s.offset, err = s.read()
	if err != nil {
		return nil, err
	}
	s.prevData = s.offset

	return s, nil
}

func (s *As5510) read() (int, error) {
	// reg 0x00 holds D7..D0 (low byte), reg 0x01 holds D9,D8 in bits [1:0]
	data, err := s.i2cChannel.Read([]byte{0x00, 0x01}, []int{1, 1})
	if err != nil {
		return -1, err
	}

	value := ((int(data[1]) & 0x03) << 8) | int(data[0])

	return value, nil
}

func (s *As5510) Read() (float64, error) {
	// read from i2c
	data, err := s.read()
	if err != nil {
		return s.prevValue, err
	}

	// calculate reset values
	diff := float64(data - s.prevData)
	change := float64(s.maxValue) * (1 - s.resetThresh)
	if diff < 0 && math.Abs(diff) > change {
		s.resetCount++
	}
	if diff > 0 && math.Abs(diff) > change {
		s.resetCount--
	}

	s.prevData = data
	output := data - s.offset + s.resetCount*(s.maxValue+1)

	position := -float64(output) * s.stepToMm
	s.prevValue = position

	return position, nil
}

func (s *As5510) Raw() (int, error) {
	return s.read()
}

// SetStepToMm sets the conversion factor from raw counts to millimeters.
// This is calibration-dependent (magnet geometry and the sensitivity register),
// so it must be set per setup rather than assumed.
func (s *As5510) SetStepToMm(stepToMm float64) {
	s.stepToMm = stepToMm
}

func (s *As5510) Close() error {
	err := s.i2cChannel.Close()
	if err != nil {
		return err
	}

	return nil
}
