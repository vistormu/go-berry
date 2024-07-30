package commands

import (
    "time"
    "fmt"
    "os"
    "math"
    "os/signal"
    "syscall"
    
    "goraspio/loadcell"
    "goraspio/motor"
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

    // variables
    programStartTime := time.Now()
    loadRef := 0.0
    refReached := false

    fmt.Println("[release] releasing wire")

    for {
    select {
    case <- quit:
        fmt.Println("\n[calibrate] stopping program")
        return nil

    case <- ticker.C:
        if time.Since(programStartTime).Seconds() > 40 && !refReached {
            fmt.Println()
            return fmt.Errorf("[calibrate] timeout error")
        }
        if time.Since(programStartTime).Seconds() > 10 && refReached {
            fmt.Println()
            return nil
        }

        // read
        _, load, err := lc.Read()
        if err != nil {
            return err
        }

        fmt.Printf("\rload: %.4f", load)

        // motor
        _, err = motor.Write(loadRef-load, 0.01)
        if err != nil {
            return err
        }

        if  math.Abs(loadRef-load) < 0.01 {
            refReached = true
        }
    }
    }
}
