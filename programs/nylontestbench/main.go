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
    phi := -math.Pi/2
    signals := []refgen.Signal{
        // ===
        // SIN
        // ===
        // 00-20%
        refgen.NewSine(20.0/2, 0.02, phi, 20.0/2), // 20mHz
        refgen.NewSine(20.0/2, 0.03, phi, 20.0/2), // 30mHz
        refgen.NewSine(20.0/2, 0.04, phi, 20.0/2), // 40mHz

        // 00-15%
        refgen.NewSine(15.0/2, 0.02, phi, 15.0/2), // 20mHz
        refgen.NewSine(15.0/2, 0.03, phi, 15.0/2), // 30mHz
        refgen.NewSine(15.0/2, 0.04, phi, 15.0/2), // 40mHz

        // 00-10%
        refgen.NewSine(10.0/2, 0.02, phi, 10.0/2), // 20mHz
        refgen.NewSine(10.0/2, 0.03, phi, 10.0/2), // 30mHz
        refgen.NewSine(10.0/2, 0.04, phi, 10.0/2), // 40mHz

        // 00-05%
        refgen.NewSine(5.0/2, 0.02, phi, 5.0/2), // 20mHz
        refgen.NewSine(5.0/2, 0.03, phi, 5.0/2), // 30mHz
        refgen.NewSine(5.0/2, 0.04, phi, 5.0/2), // 40mHz

        // 05-10%
        refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+5.0), // 20mHz
        refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+5.0), // 30mHz
        refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+5.0), // 40mHz

        // 05-15%
        refgen.NewSine(10.0/2, 0.02, phi, 10.0/2+5.0), // 20mHz
        refgen.NewSine(10.0/2, 0.03, phi, 10.0/2+5.0), // 30mHz
        refgen.NewSine(10.0/2, 0.04, phi, 10.0/2+5.0), // 40mHz

        // 05-20%
        refgen.NewSine(15.0/2, 0.02, phi, 15.0/2+5.0), // 20mHz
        refgen.NewSine(15.0/2, 0.03, phi, 15.0/2+5.0), // 30mHz
        refgen.NewSine(15.0/2, 0.04, phi, 15.0/2+5.0), // 40mHz

        // 10-15%
        refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+10.0), // 20mHz
        refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+10.0), // 30mHz
        refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+10.0), // 40mHz

        // 10-20%
        refgen.NewSine(10.0/2, 0.02, phi, 10.0/2+10.0), // 20mHz
        refgen.NewSine(10.0/2, 0.03, phi, 10.0/2+10.0), // 30mHz
        refgen.NewSine(10.0/2, 0.04, phi, 10.0/2+10.0), // 40mHz

        // 15-20%
        refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+15.0), // 20mHz
        refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+15.0), // 30mHz
        refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+15.0), // 40mHz

        // ===
        // TRI
        // ===
        // 00-20%
        refgen.NewTriangular(20.0/2, 0.02, phi, 20.0/2), // 20mHz
        refgen.NewTriangular(20.0/2, 0.03, phi, 20.0/2), // 30mHz
        refgen.NewTriangular(20.0/2, 0.04, phi, 20.0/2), // 40mHz

        // 00-15%
        refgen.NewTriangular(15.0/2, 0.02, phi, 15.0/2), // 20mHz
        refgen.NewTriangular(15.0/2, 0.03, phi, 15.0/2), // 30mHz
        refgen.NewTriangular(15.0/2, 0.04, phi, 15.0/2), // 40mHz

        // 00-10%
        refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2), // 20mHz
        refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2), // 30mHz
        refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2), // 40mHz

        // 00-05%
        refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2), // 20mHz
        refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2), // 30mHz
        refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2), // 40mHz

        // 05-10%
        refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+5.0), // 20mHz
        refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+5.0), // 30mHz
        refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+5.0), // 40mHz

        // 05-15%
        refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2+5.0), // 20mHz
        refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2+5.0), // 30mHz
        refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2+5.0), // 40mHz

        // 05-20%
        refgen.NewTriangular(15.0/2, 0.02, phi, 15.0/2+5.0), // 20mHz
        refgen.NewTriangular(15.0/2, 0.03, phi, 15.0/2+5.0), // 30mHz
        refgen.NewTriangular(15.0/2, 0.04, phi, 15.0/2+5.0), // 40mHz

        // 10-15%
        refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+10.0), // 20mHz
        refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+10.0), // 30mHz
        refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+10.0), // 40mHz

        // 10-20%
        refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2+10.0), // 20mHz
        refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2+10.0), // 30mHz
        refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2+10.0), // 40mHz

        // 15-20%
        refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+15.0), // 20mHz
        refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+15.0), // 30mHz
        refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+15.0), // 40mHz
    }

    exeInfo := ExeInfo{
        exeTime: 5*60,
        dt: 0.01,
    }

    sensorLength := 180.0 // mm
    var loadRef float32 = 1 // N

    for i, s := range signals {
        // calibration
        release()
        calibrated := calibrate(loadRef)
        if !calibrated {
            fmt.Println("Error calibrating")
            break
        }
        
        // experiment
        finalized := exe([]refgen.Signal{s}, exeInfo, sensorLength)
        if !finalized {
            break
        }

        // wait for next experiment
        if i != len(signals) - 1 {
            time.Sleep(time.Second*10)
        } 
    }
}

func exe(signals []refgen.Signal, exeInfo ExeInfo, sensorLength float64) bool {
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
    var voltageInit float32 = 0.0
    initialVoltageSet := false
    var voltageRel float32 = 0.0
    var voltageFiltRel float32 = 0.0

    // =========
    // MAIN LOOP
    // =========
    programStartTime := time.Now()
    timeFromStart := 0.0

    for range int(float64(exeInfo.exeTime)/exeInfo.dt) {
    select {
    case <- quit:
        fmt.Println("\n\nKeyboard interrupt")
        return false
    case <-ticker.C:
        loopStartTime := time.Now()

        // READ
        voltage, voltageFilt, err := vs.Read() // V
        if err != nil {
            panic(err)
        }
        if !initialVoltageSet {
            voltageInit = voltage
            initialVoltageSet = true
        }
        voltageRel = (voltage-voltageInit) / voltageInit
        voltageFiltRel = (voltageFilt-voltageInit) / voltageInit
        
        position, err := hs.Read() // mm
        if err != nil {
            position = prevPositionValue
        }
        if position != -1 {
            prevPositionValue = position
        }
        strain := position / sensorLength // strain

        load, loadFilt, err := lc.Read() // N
        if err != nil {
            panic(err)
        }

        // REFERENCE
        ref := rg.Compute(time.Since(programStartTime).Seconds()) // strain

        // ACTUATE
        _, err = motor.Write(ref-strain, exeInfo.dt)
        if err != nil {
            panic(err)
        }

        // SEND
        data["time"] = time.Since(programStartTime).Seconds()

        data["load"] = load
        data["load_filt"] = loadFilt

        data["voltage"] = voltage
        data["voltage_filt"] = voltageFilt
        data["voltage_rel"] = voltageRel
        data["voltage_filt_rel"] = voltageFiltRel

        data["reference"] = ref
        data["position"] = position
        data["strain"] = strain

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
    fmt.Println("\n\nExperiment finalized")

    return true
}


func release() {
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

    programStartTime := time.Now()

    fmt.Println("Releasing wire")

    for {
        <-ticker.C

        if time.Since(programStartTime).Seconds() > 5 {
            return
        }

        // read
        load, _, err := lc.Read()
        if err != nil {
            panic(err)
        }
        fmt.Printf("\r%.4f", load)

        if math.Abs(float64(load)) < 0.01 {
            return
        }

        // motor
        err = motor.WriteRaw(100, digitalio.High)
    }
}


func calibrate(loadRef float32) bool {
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
    ticker := time.NewTicker(time.Millisecond*100)
    defer ticker.Stop()

    programStartTime := time.Now()

    fmt.Println("Calibrating wire")

    for {
        <-ticker.C

        if time.Since(programStartTime).Seconds() > 40 {
            return false
        }

        // read
        load, _, err := lc.Read()
        if err != nil {
            panic(err)
        }
        fmt.Printf("\r%.4f", load)

        if load >= loadRef {
            fmt.Println()
            return true
        }

        // motor
        err = motor.WriteRaw(100, digitalio.Low)
    }
}
