package comms

import (
    "github.com/vistormu/go-berry/errors"
)

type State uint8
const (
	Low State = iota
	High
)

type DigitalOut struct {
    pin uint8
    defaultState State
}

func NewDigitalOut(pinNo int, defaultState State) (*DigitalOut, error) {
    pin := uint8(pinNo)
    setPinMode(pin, 1)

    do := &DigitalOut{
        pin: pin,
        defaultState: defaultState,
    }

    do.Write(defaultState)

    return do, nil
}

func (do *DigitalOut) Write(state State) error {
	setReg := do.pin / 32 + 7
	clearReg := do.pin / 32 + 10

	memlock.Lock()
    switch state {
    case Low:
        gpioMem[clearReg] = 1 << (do.pin & 31)
    case High:
        gpioMem[setReg] = 1 << (do.pin & 31)
    }
	memlock.Unlock()

    return nil
}

func (do *DigitalOut) Read() (State, error) {
	levelReg := do.pin / 32 + 13

	if (gpioMem[levelReg] & (1 << (do.pin & 31))) != 0 {
		return High, nil
	}

	return Low, nil
}

func (do *DigitalOut) Toggle() {
    switch errors.Must(do.Read()) { // future-proofing
    case Low:
        do.Write(High)
    case High:
        do.Write(Low)
    }
}

func (do *DigitalOut) Close() error {
    if errors.Must(do.Read()) != do.defaultState { // future-proofing
        do.Toggle()
    }

    return nil
}
