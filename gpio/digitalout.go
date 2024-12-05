package gpio

type DigitalOut struct {
    pin uint8
    defaultState State
}

func NewDigitalOut(pinNo int, defaultState State) DigitalOut {
    pin := uint8(pinNo)

    // set pin to output mode
	fselReg := pin / 10
	shift := (pin % 10) * 3
	f := uint32(1)
    pinMask := uint32(7)

	memlock.Lock()
	gpioMem[fselReg] = (gpioMem[fselReg] &^ (pinMask << shift)) | (f << shift)
    memlock.Unlock()

    // digitalout
    do := DigitalOut{
        pin: pin,
        defaultState: defaultState,
    }

    // Write default state
    do.Write(defaultState)

    return do
}

func (do *DigitalOut) Write(state State) {
	// Set register, 7 / 8 depending on bank
	// Clear register, 10 / 11 depending on bank
	setReg := do.pin / 32 + 7
	clearReg := do.pin / 32 + 10

	memlock.Lock()

    switch state {
    case Low:
        gpioMem[clearReg] = 1 << (do.pin & 31)
    case High:
        gpioMem[setReg] = 1 << (do.pin & 31)
    }

	memlock.Unlock() // not deferring saves ~600ns
}

func (do *DigitalOut) Read() State {
	// Input level register offset (13 / 14 depending on bank)
	levelReg := do.pin / 32 + 13

	if (gpioMem[levelReg] & (1 << (do.pin & 31))) != 0 {
		return High
	}

	return Low
}

func (do *DigitalOut) Toggle() {
    switch do.Read() {
    case Low:
        do.Write(High)
    case High:
        do.Write(Low)
    }
}

func (do *DigitalOut) Close() {
    if do.Read() != do.defaultState {
        do.Toggle()
    }
}
