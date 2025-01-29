package comms

import (
	"os"
	"sync"
	"syscall"

    "github.com/vistormu/go-berry/errors"
)

const (
	bcm2835Base = 0x20000000
	gpioOffset  = 0x200000
	clkOffset   = 0x101000
	pwmOffset   = 0x20C000
	spiOffset   = 0x204000
	intrOffset  = 0x00B000

	memLength = 4096

	GPPUPPDN0 = 57 // Pin pull-up/down for pins 15:0
	GPPUPPDN1 = 58 // Pin pull-up/down for pins 31:16
	GPPUPPDN2 = 59 // Pin pull-up/down for pins 47:32
	GPPUPPDN3 = 60 // Pin pull-up/down for pins 57:48
)

var (
	gpioBase int64
	clkBase  int64
	pwmBase  int64
	spiBase  int64
	intrBase int64

	irqsBackup uint64
)

var (
	memlock  sync.Mutex
	gpioMem  []uint32
	clkMem   []uint32
	pwmMem   []uint32
	spiMem   []uint32
	intrMem  []uint32
	gpioMem8 []uint8
	clkMem8  []uint8
	pwmMem8  []uint8
	spiMem8  []uint8
	intrMem8 []uint8
)

func init() {
	base := getBase()
	gpioBase = base + gpioOffset
	clkBase = base + clkOffset
	pwmBase = base + pwmOffset
	spiBase = base + spiOffset
	intrBase = base + intrOffset

    err := open()
    if err != nil {
        panic(err) // not the best but works right now
    }
}

func open() error {
	var file *os.File

    file, err := os.OpenFile("/dev/mem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	if os.IsPermission(err) { // try gpiomem otherwise (some extra functions like clock and pwm setting wont work)
		file, err = os.OpenFile("/dev/gpiomem", os.O_RDWR|os.O_SYNC, os.ModePerm)
	}
	if err != nil {
		return errors.New(errors.GPIO_INIT, "mem", err.Error())
	}
	defer file.Close()

	memlock.Lock()
	defer memlock.Unlock()

	gpioMem, gpioMem8, err = memMap(file.Fd(), gpioBase)
	if err != nil {
		return errors.New(errors.GPIO_INIT, "base", err.Error())
	}

	clkMem, clkMem8, err = memMap(file.Fd(), clkBase)
	if err != nil {
		return errors.New(errors.GPIO_INIT, "clk", err.Error())
	}

	pwmMem, pwmMem8, err = memMap(file.Fd(), pwmBase)
	if err != nil {
		return errors.New(errors.GPIO_INIT, "pwm", err.Error())
	}

	spiMem, spiMem8, err = memMap(file.Fd(), spiBase)
	if err != nil {
		return errors.New(errors.GPIO_INIT, "spi", err.Error())
	}

	intrMem, intrMem8, err = memMap(file.Fd(), intrBase)
	if err != nil {
		return errors.New(errors.GPIO_INIT, "intr", err.Error())
	}

	backupIRQs() // back up enabled IRQs, to restore it later

	return nil
}

func Close() error {
	enableIRQs(irqsBackup) // Return IRQs to state where it was before - just to be nice

	memlock.Lock()
	defer memlock.Unlock()
	for _, mem8 := range [][]uint8{gpioMem8, clkMem8, pwmMem8, spiMem8, intrMem8} {
		if err := syscall.Munmap(mem8); err != nil {
			return errors.New(errors.GPIO_CLOSE, err.Error())
		}
	}

	return nil
}
