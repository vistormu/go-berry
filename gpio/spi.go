package gpio

import (
    "github.com/vistormu/goraspio/errors"
)

const (
	csReg = 0
	fifoReg = 1 // TX/RX FIFO
	clkDivReg = 2
	cpol = 1 << 3
	cpha = 1 << 2
	ta = 1 << 7   // transfer active
	txd = 1 << 18 // tx fifo can accept data
	rxd = 1 << 17 // rx fifo contains data
	done = 1 << 16
)

type Spi struct {
    cs DigitalOut
    clk uint8
    mosi uint8
    miso uint8
    polarity uint8
    phase uint8
    speed int
}

func NewSpi(chipSelectPinNo, polarity, phase, speed int) (*Spi, error) {
    // begin spi
	spiMem[csReg] = 0
	if spiMem[csReg] == 0 {
		return nil, errors.New(errors.SPI_ROOT)
	}

    cs := NewDigitalOut(chipSelectPinNo, High)

    // set pins to spi
    f := uint32(4)
    clk := uint8(11)
    setPinMode(clk, f)
    mosi := uint8(10)
    setPinMode(mosi, f)
    miso := uint8(9)
    setPinMode(miso, f)

	clearSpiTxRxFifo()
	setSpiDiv(128)

    return &Spi{
        cs: cs,
        clk: clk,
        mosi: mosi,
        miso: miso,
        polarity: uint8(polarity),
        phase: uint8(phase),
        speed: speed,
    }, nil
}

func (s *Spi) setSpeed() {
	coreFreq := 250 * 1000000
	if isBCM2711() {
		coreFreq = 550 * 1000000
	}
	cdiv := uint32(coreFreq / s.speed)
	setSpiDiv(cdiv)
}

func (s *Spi) setMode() {
    switch s.polarity {
    case 0:
        spiMem[csReg] &^= cpol
    case 1:
        spiMem[csReg] |= cpol
    }

    switch s.phase {
    case 0:
        spiMem[csReg] &^= cpha
    case 1:
        spiMem[csReg] |= cpha
    }
}

func (s *Spi) exchange(data []byte) {
	clearSpiTxRxFifo()

	spiMem[csReg] |= ta

	for i := range data {
		for spiMem[csReg]&txd == 0 {}
		spiMem[fifoReg] = uint32(data[i])

		for spiMem[csReg]&rxd == 0 {}
		data[i] = byte(spiMem[fifoReg])
	}

	for spiMem[csReg]&done == 0 {}
	spiMem[csReg] &^= ta
}

func (s *Spi) Read(nBytes int) []byte {
    s.cs.Toggle()
    defer s.cs.Toggle()

    s.setSpeed()
    s.setMode()

	data := make([]byte, nBytes, nBytes)
	s.exchange(data)

    return data
}

func (s *Spi) Write(data ...byte) {
	s.exchange(append(data[:0:0], data...))
}

func (s *Spi) Close() {
    f := uint32(0) //input mode
    setPinMode(s.clk, f)
    setPinMode(s.mosi, f)
    setPinMode(s.miso, f)
    s.cs.Close()
}
