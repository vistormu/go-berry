package commands

import (
    "time"
    "fmt"
    "strconv"
    "os"
    "os/signal"
    "syscall"

    "goraspio/motor"
    "goraspio/hallsensor"
)


func reach(args []string) error {
    if len(args) != 2 {
        return fmt.Errorf("[reach] wrong number of args: expected 1 and got %d", len(args))
    }

    strainRef, err := strconv.ParseFloat(args[0], 64)
    if err != nil {
        return fmt.Errorf("[reach] error parsing strain reference\n%w", err)
    }
    sensorLength, err := strconv.ParseFloat(args[1], 64)
    if err != nil {
        return fmt.Errorf("[reach] error parsing sensor length\n%w", err)
    }

    // Motor
    pwmPinNo := 13
    freq := 500
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        return err
    }
    defer motor.Close()

    // hall sensor
    hs, err := hallsensor.NewI2C(0x40, 1)
    if err != nil {
        return err
    }
    defer hs.Close()

    // time
    ticker := time.NewTicker(time.Millisecond*10)
    defer ticker.Stop()

    // signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // variables
    programStartTime := time.Now()
    refReached := false
    prevPositionValue := 0.0

    fmt.Printf("[reach] reaching strain of %.3f %%\n", strainRef)

    for {
    select {
    case <- quit:
        fmt.Println("\n[reach] stopping program")
        return nil

    case <- ticker.C:
        if time.Since(programStartTime).Seconds() > 40 && !refReached {
            fmt.Println()
            return fmt.Errorf("[reach] timeout error")
        }
        if time.Since(programStartTime).Seconds() > 10 && refReached {
            fmt.Println()
            return nil
        }

        // read
        position, err := hs.Read() // mm
        if err != nil {
            position = prevPositionValue
        }
        if position != -1 {
            prevPositionValue = position
        }
        strain := position / sensorLength * 100 // strain (%)

        fmt.Printf("\rstrain: %.4f", strain)

        // motor
        strainErr := strainRef - strain
        _, err = motor.Write(strainErr, 0.01)
        if err != nil {
            return err
        }

        if strainErr < 0.01 {
            refReached = true
        }
    }
    }
}
