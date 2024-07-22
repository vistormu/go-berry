package commands

import (
    "time"
    "fmt"
    "strconv"
    "os"
    "os/signal"
    "syscall"

    "goraspio/digitalio"
    "goraspio/loadcell"
    "goraspio/motor"
)


func calibrate(args []string) error {
    if len(args) != 1 {
        return fmt.Errorf("[calibrate] wrong number of args: expected 1 and got %d", len(args))
    }

    loadRef, err := strconv.ParseFloat(args[0], 64)
    if err != nil {
        return fmt.Errorf("[calibrate] error parsing load reference\n%w", err)
    }

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

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    programStartTime := time.Now()

    fmt.Printf("[calibrate] calibrating wire to %.3f\n", loadRef)

    for {
    select {
    case <- quit:
        fmt.Println("\n[calibrate] stopping program")
        return nil

    case <- ticker.C:
        if time.Since(programStartTime).Seconds() > 40 {
            return fmt.Errorf("[calibrate] timeout error")
        }

        // read
        load, _, err := lc.Read()
        if err != nil {
            panic(err)
        }
        fmt.Printf("\rcurrent value %.4f", load)

        if load >= loadRef {
            fmt.Println()
            return nil
        }

        // motor
        err = motor.WriteRaw(100, digitalio.Low)
    }
    }
}
