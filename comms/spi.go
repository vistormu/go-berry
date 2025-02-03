package comms

import (
    "github.com/vistormu/go-berry/errors"
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

func init() {
	spiMem[csReg] = 0
	if spiMem[csReg] == 0 {
		panic(errors.New(errors.SPI_ROOT))
	}

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
}

type Spi struct {
    cs *DigitalOut
    polarity uint8
    phase uint8
    speed int
}

func NewSpi(chipSelectPinNo, polarity, phase, speed int) (*Spi, error) {
    return &Spi{
        cs : errors.Must(NewDigitalOut(chipSelectPinNo, High)), // future-proofing
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

func (s *Spi) Exchange(data []byte) {
    s.cs.Toggle()
    defer s.cs.Toggle()

    s.setSpeed()
    s.setMode()

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

func (s *Spi) Read(nBytes int) ([]byte, error) {
	data := make([]byte, nBytes, nBytes)
	s.Exchange(data)

    return data, nil
}

func (s *Spi) Write(data ...byte) error {
	s.Exchange(append(data[:0:0], data...))

    return nil
}

func (s *Spi) Close() error {
    f := uint32(0) //input mode
    clk := uint8(11)
    setPinMode(clk, f)
    mosi := uint8(10)
    setPinMode(mosi, f)
    miso := uint8(9)
    setPinMode(miso, f)

    err := s.cs.Close()
    if err != nil {
        return err
    }

    return nil
}
