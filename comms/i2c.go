package comms

import (
	"fmt"
	"os"

    "github.com/vistormu/go-berry/errors"
)

 // #include <linux/i2c-dev.h>
import "C"

type I2C struct {
	addr uint8
	bus  int
	rc   *os.File
}

func NewI2C(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, errors.New(errors.I2C_OPEN, err.Error())
	}

    err = ioctl(f.Fd(), C.I2C_SLAVE, uintptr(addr))
    if err != nil {
        return nil, errors.New(errors.I2C_OPEN, err.Error())
    }

    return &I2C{
        rc: f, 
        bus: bus, 
        addr: addr,
    }, nil
}

func (i *I2C) Read(registers []byte, nBytes []int) ([]byte, error) {
    if len(registers) > 0 && isConsecutive(registers) {
        startRegister := registers[0]
        totalBytes := sum(nBytes)

        return i.readBulk(startRegister, totalBytes)
    }

    result := make([]byte, sum(nBytes))
    offset := 0

    for j, reg := range registers {
        n := nBytes[j]

        if _, err := i.rc.Write([]byte{reg}); err != nil {
            return nil, errors.New(errors.I2C_READ, reg, err.Error())
        }

        buf := make([]byte, n)
        if _, err := i.rc.Read(buf); err != nil {
            return nil, errors.New(errors.I2C_READ, reg, err.Error())
        }

        copy(result[offset:], buf)
        offset += n
    }

    return result, nil
}

func (i *I2C) Write(reg byte, value byte) error {
	buf := []byte{reg, value}
	_, err := i.rc.Write(buf)
	if err != nil {
		return errors.New(errors.I2C_WRITE, reg, err.Error())
	}

	return nil
}

func (i *I2C) Close() error {
    err := i.rc.Close()
    if err != nil {
        return errors.New(errors.I2C_CLOSE, err.Error())
    }

	return nil
}

func (i *I2C) readBulk(startRegister byte, nBytes int) ([]byte, error) {
    if _, err := i.rc.Write([]byte{startRegister}); err != nil {
        return nil, errors.New(errors.I2C_READ, startRegister, err.Error())
    }

    result := make([]byte, nBytes)
    if _, err := i.rc.Read(result); err != nil {
        return nil, errors.New(errors.I2C_READ, startRegister, err.Error())
    }

    return result, nil
}

func isConsecutive(registers []byte) bool {
    for i := 1; i < len(registers); i++ {
        if registers[i] != registers[i-1]+1 {
            return false
        }
    }
    return true
}

func sum(values []int) int {
    result := 0
    for _, value := range values {
        result += value
    }

    return result
}
