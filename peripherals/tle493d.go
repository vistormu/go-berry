package peripherals

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/vistormu/go-berry/comms"
)

// ==================
// registers and maps
// ==================
type Reg struct {
	reg   byte
	mask  byte
	shift uint
}

var readRegs = map[string]Reg{
	"bxTop":         {0x00, 0xFF, 0},
	"bxBottom":      {0x04, 0xF0, 4},
	"byTop":         {0x01, 0xFF, 0},
	"byBottom":      {0x04, 0x0F, 0},
	"bzTop":         {0x02, 0xFF, 0},
	"bzBottom":      {0x05, 0x0F, 0},
	"temp1":         {0x03, 0xF0, 4},
	"temp2":         {0x06, 0xFF, 0},
	"framCounter":   {0x03, 0x0C, 2},
	"channel":       {0x03, 0x03, 0},
	"powerDownFlag": {0x05, 0x10, 4},
	"res1":          {0x07, 0x18, 3},
	"res2":          {0x08, 0xFF, 0},
	"res3":          {0x09, 0x1F, 0},
	"testMode":      {0x05, 0x40, 6},
	"parityFuse":    {0x05, 0x20, 5},
}

var writeRegs = map[string]Reg{
	"parity":      {0x01, 0x80, 7},
	"addr":        {0x01, 0x60, 5},
	"int":         {0x01, 0x04, 2},
	"fast":        {0x01, 0x02, 1},
	"lowPower":    {0x01, 0x01, 0},
	"tempDisable": {0x03, 0x80, 7},
	"lpPeriod":    {0x03, 0x40, 6},
	"parityTest":  {0x03, 0x20, 5},
	"powerDown":   {0x03, 0x20, 5},
	"res1":        {0x01, 0x18, 3},
	"res2":        {0x02, 0xFF, 0},
	"res3":        {0x03, 0x1F, 0},
}

type Tle493d struct {
	i2c *comms.I2C

	readBuffer []byte
	readRegs   []byte
	readNBytes []int

	prevAngle float64
}

func NewTle493d(address uint8, line int) (*Tle493d, error) {
	i2c, err := comms.NewI2C(0x22, line)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the first address\n%v", err)
	}

	// writing 0x10 to the 0x11 registry sets it to an special mode?
	err = i2c.Write(0x11, 0x13)
	if err != nil {
		fmt.Println("error writing 0x13 to 0x11")
	}
	i2c.Close()

	// the i2c address changes to 0x35
	i2c, err = comms.NewI2C(address, line)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the sensor\n%v", err)
	}

	// reading buffer variables
	bytes := make([]int, 10)
	regs := make([]byte, 10)
	for i := range 10 {
		regs[i] = byte(i)
		bytes[i] = 1
	}

	s := &Tle493d{
		i2c:        i2c,
		readBuffer: make([]byte, 10),
		readRegs:   regs,
		readNBytes: bytes,
		prevAngle:  0,
	}

	// initialize read buffer
	err = s.update()
	if err != nil {
		return nil, fmt.Errorf("error initilizing read buffer\n%v", err)
	}

	buffer := make([]byte, 4)

	// setup write registers
	for _, k := range []string{"res1", "res2", "res3"} {
		value := s.get(k)
		s.write(buffer, k, value)
	}

	// write the correct I2C address
	s.write(buffer, "addr", 0)

	// set "master controlled mode" (take measurement on every read)
	s.write(buffer, "parity", 1)
	s.write(buffer, "fast", 1)
	s.write(buffer, "lowPower", 1)

	// write to the I2C
	for i, b := range buffer {
		err := s.i2c.Write(byte(i), b)
		if err != nil {
			return nil, fmt.Errorf("error writing to register %d\n%v", i, err)
		}
	}

	return s, nil
}

func (s *Tle493d) write(buffer []byte, key string, value int) {
	r := writeRegs[key]
	current := buffer[int(r.reg)]
	current &= ^r.mask
	current |= (byte(value) << r.shift) & r.mask

	buffer[int(r.reg)] = current
}

func (s *Tle493d) update() error {
	data, err := s.i2c.Read(s.readRegs, s.readNBytes)
	if err != nil {
		return err
	}

	s.readBuffer = data

	// fmt.Println(data)

	return nil
}

func (s *Tle493d) get(key string) int {
	r := readRegs[key]
	value := s.readBuffer[int(r.reg)]

	return int((value & r.mask) >> r.shift)
}

func combine(top, bottom int) float64 {
	buffer := []byte{
		byte(top),
		byte((bottom << 4) & 0xFF),
	}

	value := int16(binary.BigEndian.Uint16(buffer))
	value = value >> 4

	return float64(value) * 98.0
}

func (s *Tle493d) read() (float64, float64, error) {
	err := s.update()
	if err != nil {
		return 0.0, 0.0, err
	}

	bxTop := s.get("bxTop")
	bxBot := s.get("bxBot")
	byTop := s.get("byTop")
	byBot := s.get("byBot")

	bx := combine(bxTop, bxBot)
	by := combine(byTop, byBot)

	return bx, by, nil
}

func (s *Tle493d) Read() (float64, float64, error) {
	err := s.update()
	if err != nil {
		return 0.0, 0.0, err
	}

	bxTop := s.get("bxTop")
	bxBot := s.get("bxBot")
	byTop := s.get("byTop")
	byBot := s.get("byBot")

	// fmt.Print(bxTop, bxTop, "\n")

	bx := combine(bxTop, bxBot)
	by := combine(byTop, byBot)

	return bx, by, nil
}

func (s *Tle493d) Position() (float64, error) {
	// UNDER DEVELOPMENT
	bx, _, err := s.read()
	if err != nil {
		return s.prevAngle, err
	}

	return bx, nil
}

func (s *Tle493d) Angle() (float64, error) {
	bx, by, err := s.read()
	if err != nil {
		return s.prevAngle, err
	}

	angle := math.Atan2(by, bx)
	angle *= 180 / math.Pi
	// angle = math.Mod(angle + 360, 360)

	s.prevAngle = angle

	fmt.Println(angle)

	return angle, nil
}

func (s *Tle493d) Close() error {
	err := s.i2c.Close()
	if err != nil {
		return err
	}

	return nil
}
