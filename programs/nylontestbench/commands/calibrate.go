package commands

import (
    "time"
    "fmt"
    "strconv"
    "os"
    "os/signal"
    "syscall"

    "goraspio/loadcell"
    "goraspio/motor"
    "goraspio/voltagesensor"
    "goraspio/utils"
)


func calibrate(args []string) (float64, error) {
    if len(args) != 1 {
        return 0.0, fmt.Errorf("[calibrate] wrong number of args: expected 1 and got %d", len(args))
    }

    loadRef, err := strconv.ParseFloat(args[0], 64)
    if err != nil {
        return 0.0, fmt.Errorf("[calibrate] error parsing load reference\n%w", err)
    }

    // Motor
    pwmPinNo := 13
    freq := 200
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        return 0.0, err
    }
    defer motor.Close()

    // Load cell
    lc, err := loadcell.New(24)
    if err != nil {
        return 0.0, err
    }
    defer lc.Close()

    // voltage sensor
    vs, err := voltagesensor.New(5, 25, 0.0)
    if err != nil {
        return 0.0, err
    }
    defer vs.Close()

    // time
    ticker := time.NewTicker(time.Millisecond*10)
    defer ticker.Stop()

    // signals
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

    // variables
    programStartTime := time.Now()
    var timeReached time.Time
    refReached := false
    voltageValues := make([]float64, 0)

    fmt.Printf("[calibrate] calibrating wire to %.3f\n", loadRef)

    for {
    select {
    case <- quit:
        fmt.Println("\n[calibrate] stopping program")
        return 0.0, nil

    case <- ticker.C:
        if time.Since(programStartTime).Seconds() > 40 && !refReached {
            fmt.Println()
            return 0.0, fmt.Errorf("[calibrate] timeout error")
        }
        if !timeReached.IsZero() {
            if time.Since(timeReached).Seconds() > 10 {
                fmt.Println()
                return utils.Mean(voltageValues), nil
            }
        }

        // read
        load, loadFilt, err := lc.Read()
        if err != nil {
            return 0.0, err
        }
        vr, err := vs.Read()
        if err != nil {
            return 0.0, err
        }

        fmt.Printf("\rload: %.4f | loadFilt: %.4f | voltage: %.4f", load, loadFilt, vr.VoltageFilt)

        // motor
        loadErr := loadRef - loadFilt
        _, err = motor.Write(loadErr, 0.01)
        if err != nil {
            return 0.0, err
        }

        if loadErr < 0.01  && !refReached {
            refReached = true
            timeReached = time.Now()
        }

        if refReached {
            voltageValues = append(voltageValues, vr.VoltageFilt)
        }
    }
    }
}
