package main

import (
	"fmt"
	"goraspio/client"
	"goraspio/hallsensor"
	"goraspio/motor"
	"goraspio/refgen"
	"goraspio/voltagesensor"
	"os"
	"os/signal"
	"syscall"
	"time"
    "math"
)

type ExeInfo struct {
    exeTime int
    dt float64
}

func main() {
    amp := 20.0/2
    freq := 0.04
    phi := -math.Pi/2
    offset := amp
    sine := refgen.NewSine(amp, freq, phi, offset)

    signals := make([]refgen.Signal, 1)
    signals[0] = sine

    exeInfo := ExeInfo{
        exeTime: 25,
        dt: 0.01,
    }

    exe(signals, exeInfo)
}

func exe(signals []refgen.Signal, exeInfo ExeInfo) {
    // ==========
    // COMPONENTS
    // ==========
    // Motor
    pwmPinNo := 13
    freq := 10_000
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        panic(err)
    }
    defer motor.Close()
    fmt.Println("motor connected successfully")

    // Voltage Sensor
    var vRef float32 = 5.0
    voltageSensorhipSelectNo := 25
    vs, err := voltagesensor.New(vRef, voltageSensorhipSelectNo)
    if err != nil {
        panic(err)
    }
    defer vs.Close()
    fmt.Println("voltage sensor connected successfully")

    // Hall Sensor
    hs, err := hallsensor.NewI2C(0x40, 1)
    if err != nil {
        panic(err)
    }
    defer hs.Close()
    fmt.Println("hall sensor connected successfully")

    // Client
    ip := "10.118.90.193"
    port := 8080
    c, err := client.New(ip, port)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    fmt.Println("client connected successfully")

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
            fmt.Println("WARNING: I2C failed reading")
            // panic(err)
        }

        // REFERENCE
        ref := rg.Compute(time.Since(programStartTime).Seconds())

        // ERROR
        posError := ref - position

        // ACTUATE
        pwmValue, err := motor.Write(posError)
        if err != nil {
            panic(err)
        }

        // SEND
        data["time"] = time.Since(programStartTime).Seconds()

        data["control"] = pwmValue

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
