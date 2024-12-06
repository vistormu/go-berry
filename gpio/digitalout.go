package gpio

type State uint8
const (
	Low State = iota
	High
)

type DigitalOut struct {
    pin uint8
    defaultState State
}

func NewDigitalOut(pinNo int, defaultState State) DigitalOut {
    pin := uint8(pinNo)
    setPinMode(pin, 1)

    do := DigitalOut{
        pin: pin,
        defaultState: defaultState,
    }

    do.Write(defaultState)

    return do
}

func (do DigitalOut) Write(state State) {
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
}

func (do DigitalOut) Read() State {
	levelReg := do.pin / 32 + 13

	if (gpioMem[levelReg] & (1 << (do.pin & 31))) != 0 {
		return High
	}

	return Low
}

func (do DigitalOut) Toggle() {
    switch do.Read() {
    case Low:
        do.Write(High)
    case High:
        do.Write(Low)
    }
}

func (do DigitalOut) Close() {
    if do.Read() != do.defaultState {
        do.Toggle()
    }
}
