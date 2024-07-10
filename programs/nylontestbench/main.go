package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
    "math"
    
    "goraspio/digitalio"
	"goraspio/client"
	"goraspio/hallsensor"
	"goraspio/motor"
	"goraspio/refgen"
	"goraspio/voltagesensor"
    "goraspio/loadcell"
)

type ExeInfo struct {
    exeTime int
    dt float64
}

func main() {
    // phi := -math.Pi/2
    signals := []refgen.Signal{
        // refgen.NewSine(20.0/2, 0.02, phi, 20.0/2), // sin 20mHz 0-20mm
        // refgen.NewSine(20.0/2, 0.03, phi, 20.0/2), // sin 30 0-20
        // refgen.NewSine(20.0/2, 0.04, phi, 20.0/2), // sin 40 0-20
        // refgen.NewSine(10.0/2, 0.02, phi, 10.0/2), // sin 20 0-10
        // refgen.NewSine(10.0/2, 0.03, phi, 10.0/2), // sin 30 0-10
        // refgen.NewSine(10.0/2, 0.04, phi, 10.0/2), // sin 40 0-10
        // refgen.NewSine(5.0/2, 0.02, phi, 5.0/2), // sin 20 0-5
        // refgen.NewSine(5.0/2, 0.03, phi, 5.0/2), // sin 30 0-5
        // refgen.NewSine(5.0/2, 0.04, phi, 5.0/2), // sin 40 0-5
        // refgen.NewSine(15.0/2, 0.02, phi, 15.0/2), // sin 20 0-15
        // refgen.NewSine(15.0/2, 0.03, phi, 15.0/2), // sin 30 0-15
        // refgen.NewSine(15.0/2, 0.04, phi, 15.0/2), // sin 40 0-15
    }

    exeInfo := ExeInfo{
        exeTime: 5*60,
        dt: 0.01,
    }

    var loadRef float32 = 0.15

    for _, s := range signals {
        calibrate(loadRef)
        exe([]refgen.Signal{s}, exeInfo)
        time.Sleep(time.Second*10)
    }
}

func exe(signals []refgen.Signal, exeInfo ExeInfo) {
    // ==========
    // COMPONENTS
    // ==========
    // Motor
    pwmPinNo := 13
    freq := 4_500
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        panic(err)
    }
    defer motor.Close()

    // Voltage Sensor
    var vRef float32 = 5.0
    voltageSensorhipSelectNo := 25
    vs, err := voltagesensor.New(vRef, voltageSensorhipSelectNo)
    if err != nil {
        panic(err)
    }
    defer vs.Close()

    // Hall Sensor
    hs, err := hallsensor.NewI2C(0x40, 1)
    if err != nil {
        panic(err)
    }
    defer hs.Close()

    // Load cell
    lc, err := loadcell.New(24)
    if err != nil {
        panic(err)
    }
    defer lc.Close()

    // Client
    ip := "10.118.90.193"
    port := 8080
    c, err := client.New(ip, port)
    if err != nil {
        panic(err)
    }
    defer c.Close()

    // Reference generator
    rg := refgen.NewRefGen(signals)

    // ========
    // CHANNELS
    // ========
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    ticker := time.NewTicker(time.Duration(exeInfo.dt*float64(time.Second)))
    defer ticker.Stop()

    // =========
    // VARIABLES
    // =========
    data := make(map[string]any)
    prevPositionValue := 0.0

    // =========
    // MAIN LOOP
    // =========
    programStartTime := time.Now()
    timeFromStart := 0.0

    for range int(float64(exeInfo.exeTime)/exeInfo.dt) {
    select {
    case <- quit:
        fmt.Println("\nexiting")
        return
    case <-ticker.C:
        loopStartTime := time.Now()

        // READ
        voltage, filteredVoltage, err := vs.Read()
        if err != nil {
            panic(err)
        }
        position, err := hs.Read()
        if err != nil {
            position = prevPositionValue
        }
        if position != -1 {
            prevPositionValue = position
        }
        load, filteredLoad, err := lc.Read()
        if err != nil {
            panic(err)
        }

        // REFERENCE
        ref := rg.Compute(time.Since(programStartTime).Seconds())

        // ERROR
        posError := ref - position

        // ACTUATE
        _, err = motor.Write(posError, exeInfo.dt)
        if err != nil {
            panic(err)
        }

        // SEND
        data["time"] = time.Since(programStartTime).Seconds()

        data["load"] = load
        data["filtered_load"] = filteredLoad

        data["voltage"] = voltage
        data["filtered_voltage"] = filteredVoltage

        data["reference"] = ref
        data["position"] = position

        err = c.Send(data)
        if err != nil {
            panic(err)
        }

        // TIME
        timePerIteration := time.Since(loopStartTime).Seconds()*1000
        timeFromStart = time.Since(programStartTime).Seconds()

        // PRINT
        fmt.Printf("\rTime: %.3f ms / %.3f s | Position: %.3f mm", timePerIteration, timeFromStart, position)
    }}
    fmt.Println("\nProgram finalized")
}


func calibrate(loadRef float32) {
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
    directionSet := false
    direction := digitalio.Low

    for {
        <-ticker.C

        if time.Since(programStartTime).Seconds() > 20 {
            break
        }

        // read
        _, load, err := lc.Read()
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
