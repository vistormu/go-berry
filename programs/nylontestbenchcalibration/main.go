package main

import (
    "time"
    "fmt"

    "goraspio/digitalio"
    "goraspio/loadcell"
    "goraspio/motor"
)

func main() {
    // Motor
    pwmPinNo := 13
    freq := 500
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        panic(err)
    }
    defer motor.Close()

    // Load cell
    lc, err := loadcell.New(24)
    if err != nil {
        panic(err)
    }
    defer lc.Close()

    // time
    ticker := time.NewTicker(time.Millisecond*10)
    defer ticker.Stop()

    programStartTime := time.Now()

    // variables
    var loadRef float32 = 0.2
    directionSet := false
    direction := digitalio.Low

    for {
        <-ticker.C

        if time.Since(programStartTime).Seconds() > 20 {
            break
        }

        // read
        load, err := lc.Read()
        if err != nil {
            panic(err)
        }
        fmt.Printf("\r%.4f", load)

        // check initial direction
        if load < loadRef && !directionSet { // positive
            direction = digitalio.Low
            directionSet = true
        } else if load > loadRef && !directionSet {
            direction = digitalio.High // left
            directionSet = true
        }

        // move motor
        if load >= loadRef && direction == digitalio.Low { // load went right and surpassed the ref
            break
        }
        if load <= loadRef && direction == digitalio.High { // load went left and surpasses the ref
            break
        }

        // motor
        err = motor.WriteRaw(100, direction)
    }
}
