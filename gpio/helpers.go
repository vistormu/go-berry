package gpio

import (
    "os"
	"bytes"
	"encoding/binary"
	"errors"
	"unsafe"
	"reflect"
	"syscall"
)

// Read /proc/device-tree/soc/ranges and determine the base address.
// Use the default Raspberry Pi 1 base address if this fails.
func readBase(offset int64) (int64, error) {
	ranges, err := os.Open("/proc/device-tree/soc/ranges")
	defer ranges.Close()
	if err != nil {
		return 0, err
	}
	b := make([]byte, 4)
	n, err := ranges.ReadAt(b, offset)
	if n != 4 || err != nil {
		return 0, err
	}
	buf := bytes.NewReader(b)
	var out uint32
	err = binary.Read(buf, binary.BigEndian, &out)
	if err != nil {
		return 0, err
	}

	if out == 0 {
		return 0, errors.New("rpio: GPIO base address not found")
	}
	return int64(out), nil
}

func getBase() int64 {
	// Pi 2 & 3 GPIO base address is at offset 4
	b, err := readBase(4)
	if err == nil {
		return b
	}

	// Pi 4 GPIO base address is as offset 8
	b, err = readBase(8)
	if err == nil {
		return b
	}

	// Default to Pi 1
	return int64(bcm2835Base)
}

// The Pi 4 uses a BCM 2711, which has different register offsets and base addresses than the rest of the Pi family (so far).  This
// helper function checks if we're on a 2711 and hence a Pi 4
func isBCM2711() bool {
	return gpioMem[GPPUPPDN3] != 0x6770696f
}

func memMap(fd uintptr, base int64) (mem []uint32, mem8 []byte, err error) {
	mem8, err = syscall.Mmap(
		int(fd),
		base,
		memLength,
		syscall.PROT_READ|syscall.PROT_WRITE,
		syscall.MAP_SHARED,
	)
	if err != nil {
		return
	}
	// Convert mapped byte memory to unsafe []uint32 pointer, adjust length as needed
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&mem8))
	header.Len /= (32 / 8) // (32 bit = 4 bytes)
	header.Cap /= (32 / 8)
	mem = *(*[]uint32)(unsafe.Pointer(&header))
	return
}

func backupIRQs() {
	const irqEnable1 = 0x210 / 4
	const irqEnable2 = 0x214 / 4
	irqsBackup = uint64(intrMem[irqEnable2])<<32 | uint64(intrMem[irqEnable1])
}

func enableIRQs(irqs uint64) {
	const irqEnable1 = 0x210 / 4
	const irqEnable2 = 0x214 / 4
	intrMem[irqEnable1] = uint32(irqs)       // IRQ 0..31
	intrMem[irqEnable2] = uint32(irqs >> 32) // IRQ 32..63
}

func setPinMode(pin uint8, f uint32) {
	fselReg := pin / 10
	shift := (pin % 10) * 3
    pinMask := uint32(7)

	memlock.Lock()
	gpioMem[fselReg] = (gpioMem[fselReg] &^ (pinMask << shift)) | (f << shift)
    memlock.Unlock()
}
