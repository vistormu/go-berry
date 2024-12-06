package gpio

import (
	"errors"
)

type SpiDev int

// SPI devices.
// Only SPI0 supported for now.
const (
	Spi0 SpiDev = iota
	Spi1        // aux
	Spi2        // aux
)

const (
	csReg     = 0
	fifoReg   = 1 // TX/RX FIFO
	clkDivReg = 2
)

type Mode uint8
type Pin uint8

const (
	Input Mode = iota
	Output
	Clock
	Pwmm
	Spi
	Alt0
	Alt1
	Alt2
	Alt3
	Alt4
	Alt5
)

// Note that you should disable SPI interface in raspi-config first!
func SpiBegin(dev SpiDev) error {
	spiMem[csReg] = 0 // reset spi settings to default
	if spiMem[csReg] == 0 {
		// this should not read only zeroes after reset -> mem map failed
		return errors.New("SPI registers not mapped correctly - are you root?")
	}

	for _, pin := range getSpiPins(dev) {
		pin.Mode(Spi)
	}

	clearSpiTxRxFifo()
	setSpiDiv(128)
	return nil
}

// SpiEnd: Sets SPI pins of given device to default (Input) mode. See SpiBegin.
func SpiEnd(dev SpiDev) {
	var pins = getSpiPins(dev)
	for _, pin := range pins {
		pin.Mode(Input)
	}
}

// SpiSpeed: Set (maximal) speed [Hz] of SPI clock.
// Param speed may be as big as 125MHz in theory, but
// only values up to 31.25MHz are considered relayable.
func SpiSpeed(speed int) {
	coreFreq := 250 * 1000000
	if isBCM2711() {
		coreFreq = 550 * 1000000
	}
	cdiv := uint32(coreFreq / speed)
	setSpiDiv(cdiv)
}

// SpiChipSelect: Select chip, one of 0, 1, 2
// for selecting slave on CE0, CE1, or CE2 pin
func SpiChipSelect(chip uint8) {
	const csMask = 3 // chip select has 2 bits

	cs := uint32(chip & csMask)

	spiMem[csReg] = spiMem[csReg]&^csMask | cs
}

// SpiChipSelectPolarity: Sets polarity (0/1) of active chip select
// default active=0
func SpiChipSelectPolarity(chip uint8, polarity uint8) {
	if chip > 2 {
		return
	}
	cspol := uint32(1 << (21 + chip)) // bit 21, 22 or 23 depending on chip

	if polarity == 0 { // chip select is active low
		spiMem[csReg] &^= cspol
	} else { // chip select is active hight
		spiMem[csReg] |= cspol
	}
}

// SpiMode: Set polarity (0/1) and phase (0/1) of spi clock
// default polarity=0; phase=0
func SpiMode(polarity uint8, phase uint8) {
	const cpol = 1 << 3
	const cpha = 1 << 2

	if polarity == 0 { // Rest state of clock = low
		spiMem[csReg] &^= cpol
	} else { // Rest state of clock = high
		spiMem[csReg] |= cpol
	}

	if phase == 0 { // First SCLK transition at middle of data bit
		spiMem[csReg] &^= cpha
	} else { // First SCLK transition at beginning of data bit
		spiMem[csReg] |= cpha
	}
}

// SpiTransmit takes one or more bytes and send them to slave.
//
// Data received from slave are ignored.
// Use spread operator to send slice of bytes.
func SpiTransmit(data ...byte) {
	SpiExchange(append(data[:0:0], data...)) // clone data because it will be rewriten by received bytes
}

// SpiReceive receives n bytes from slave.
//
// Note that n zeroed bytes are send to slave as side effect.
func SpiReceive(n int) []byte {
	data := make([]byte, n, n)
	SpiExchange(data)
	return data
}

// SpiExchange: Transmit all bytes in data to slave
// and simultaneously receives bytes from slave to data.
//
// If you want to only send or only receive, use SpiTransmit/SpiReceive
func SpiExchange(data []byte) {
	const ta = 1 << 7   // transfer active
	const txd = 1 << 18 // tx fifo can accept data
	const rxd = 1 << 17 // rx fifo contains data
	const done = 1 << 16

	clearSpiTxRxFifo()

	// set TA = 1
	spiMem[csReg] |= ta

	for i := range data {
		// wait for TXD
		for spiMem[csReg]&txd == 0 {
		}
		// write bytes to SPI_FIFO
		spiMem[fifoReg] = uint32(data[i])

		// wait for RXD
		for spiMem[csReg]&rxd == 0 {
		}
		// read bytes from SPI_FIFO
		data[i] = byte(spiMem[fifoReg])
	}

	// wait for DONE
	for spiMem[csReg]&done == 0 {
	}

	// Set TA = 0
	spiMem[csReg] &^= ta
}

// set spi clock divider value
func setSpiDiv(div uint32) {
	const divMask = 1<<16 - 1 - 1 // cdiv have 16 bits and must be odd (for some reason)
	spiMem[clkDivReg] = div & divMask
}

// clear both FIFOs
func clearSpiTxRxFifo() {
	const clearTxRx = 1<<5 | 1<<4
	spiMem[csReg] |= clearTxRx
}

func getSpiPins(dev SpiDev) []Pin {
	switch dev {
	case Spi0:
		return []Pin{7, 8, 9, 10, 11}
		// ommit 35, 36, 37, 38, 39 - only one set of SPI0 can be set in Spi mode at a time
	case Spi1:
		return []Pin{16, 17, 18, 19, 20, 21}
	case Spi2:
		return []Pin{40, 41, 42, 43, 44, 45}
	default:
		return []Pin{}
	}
}

func (pin Pin) Mode(mode Mode) {
	PinMode(pin, mode)
}
func PinMode(pin Pin, mode Mode) {

	// Pin fsel register, 0 or 1 depending on bank
	fselReg := uint8(pin) / 10
	shift := (uint8(pin) % 10) * 3
	f := uint32(0)

	const in = 0   // 000
	const out = 1  // 001
	const alt0 = 4 // 100
	const alt1 = 5 // 101
	const alt2 = 6 // 110
	const alt3 = 7 // 111
	const alt4 = 3 // 011
	const alt5 = 2 // 010

	switch mode {
	case Input:
		f = in
	case Output:
		f = out
	case Clock:
		switch pin {
		case 4, 5, 6, 32, 34, 42, 43, 44:
			f = alt0
		case 20, 21:
			f = alt5
		default:
			return
		}
	case Pwmm:
		switch pin {
		case 12, 13, 40, 41, 45:
			f = alt0
		case 18, 19:
			f = alt5
		default:
			return
		}
	case Spi:
		switch pin {
		case 7, 8, 9, 10, 11: // SPI0
			f = alt0
		case 35, 36, 37, 38, 39: // SPI0
			f = alt0
		case 16, 17, 18, 19, 20, 21: // SPI1
			f = alt4
		case 40, 41, 42, 43, 44, 45: // SPI2
			f = alt4
		default:
			return
		}
	case Alt0:
		f = alt0
	case Alt1:
		f = alt1
	case Alt2:
		f = alt2
	case Alt3:
		f = alt3
	case Alt4:
		f = alt4
	case Alt5:
		f = alt5
	}

	memlock.Lock()
	defer memlock.Unlock()

	const pinMask = 7 // 111 - pinmode is 3 bits

	gpioMem[fselReg] = (gpioMem[fselReg] &^ (pinMask << shift)) | (f << shift)
}
