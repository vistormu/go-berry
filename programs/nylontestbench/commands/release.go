package commands

import (
    "time"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "math"
    
    "goraspio/loadcell"
    "goraspio/motor"
    "goraspio/digitalio"
)


func release(args []string) error {
    if len(args) != 0 {
        return fmt.Errorf("[release] wring number of args: expected 0 and got %d", len(args))
    }

    // Motor
    pwmPinNo := 13
    freq := 2_000
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

    fmt.Println("[release] releasing wire")

    for {
    select {
    case <- quit:
        fmt.Println("\n[calibrate] stopping program")
        return nil

    case <- ticker.C:
        if time.Since(programStartTime).Seconds() > 5 {
            return fmt.Errorf("[release] timeout error")
        }

        // read
        load, _, err := lc.Read()
        if err != nil {
            panic(err)
        }
        fmt.Printf("\rcurrent value: %.4f", load)

        if math.Abs(load) < 0.01 {
            fmt.Println()
            return nil
        }

        // motor
        err = motor.WriteRaw(100, digitalio.High)
    }
    }
}
