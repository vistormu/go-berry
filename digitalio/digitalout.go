package digitalio


import (
    "github.com/stianeikeland/go-rpio/v4"
)

type PinState int
const (
    Low PinState = iota
    High
)


type DigitalOut struct {
    pin rpio.Pin
}

func NewDigitalOut(pinNo int, defaultState PinState) DigitalOut {
    pin := rpio.Pin(pinNo)

    pin.Output()
    
    if defaultState == Low {
        pin.Low()
    } else {
        pin.High()
    }

    return DigitalOut{pin}
}

func (do DigitalOut) Write(state PinState) {
    if state == Low {
        do.pin.Low()
    } else {
        do.pin.High()
    }
}

func (do DigitalOut) Toggle() PinState {
    state := do.pin.Read()

    if state == rpio.Low {
        do.pin.High()
        return High
    }

    do.pin.Low()
    return Low
}

func (do DigitalOut) Read() PinState {
    state := do.pin.Read()
    if state == rpio.Low {
        return High
    }
    return Low
}

func (do DigitalOut) Close() {
    do.pin.Low()
}
