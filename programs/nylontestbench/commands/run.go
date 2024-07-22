package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
    "math"

    // "gopkg.in/yaml.v3"
    
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

func run(args []string) error {
    phi := -math.Pi/2
    signals := []refgen.Signal{
        // ===
        // SIN
        // ===
        // 00-20%
        // refgen.NewSine(20.0/2, 0.02, phi, 20.0/2), // 20mHz
        // refgen.NewSine(20.0/2, 0.03, phi, 20.0/2), // 30mHz
        // refgen.NewSine(20.0/2, 0.04, phi, 20.0/2), // 40mHz

        // 00-15%
        // refgen.NewSine(15.0/2, 0.02, phi, 15.0/2), // 20mHz
        // refgen.NewSine(15.0/2, 0.03, phi, 15.0/2), // 30mHz
        // refgen.NewSine(15.0/2, 0.04, phi, 15.0/2), // 40mHz

        // 00-10%
        // refgen.NewSine(10.0/2, 0.02, phi, 10.0/2), // 20mHz
        // refgen.NewSine(10.0/2, 0.03, phi, 10.0/2), // 30mHz
        refgen.NewSine(10.0/2, 0.04, phi, 10.0/2), // 40mHz

        // 00-05%
        // refgen.NewSine(5.0/2, 0.02, phi, 5.0/2), // 20mHz
        // refgen.NewSine(5.0/2, 0.03, phi, 5.0/2), // 30mHz
        // refgen.NewSine(5.0/2, 0.04, phi, 5.0/2), // 40mHz

        // 05-10%
        // refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+5.0), // 20mHz
        // refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+5.0), // 30mHz
        // refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+5.0), // 40mHz

        // 05-15%
        // refgen.NewSine(10.0/2, 0.02, phi, 10.0/2+5.0), // 20mHz
        // refgen.NewSine(10.0/2, 0.03, phi, 10.0/2+5.0), // 30mHz
        // refgen.NewSine(10.0/2, 0.04, phi, 10.0/2+5.0), // 40mHz

        // 05-20%
        // refgen.NewSine(15.0/2, 0.02, phi, 15.0/2+5.0), // 20mHz
        // refgen.NewSine(15.0/2, 0.03, phi, 15.0/2+5.0), // 30mHz
        // refgen.NewSine(15.0/2, 0.04, phi, 15.0/2+5.0), // 40mHz

        // 10-15%
        // refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+10.0), // 20mHz
        // refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+10.0), // 30mHz
        // refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+10.0), // 40mHz

        // 10-20%
        // refgen.NewSine(10.0/2, 0.02, phi, 10.0/2+10.0), // 20mHz
        // refgen.NewSine(10.0/2, 0.03, phi, 10.0/2+10.0), // 30mHz
        // refgen.NewSine(10.0/2, 0.04, phi, 10.0/2+10.0), // 40mHz

        // 15-20%
        // refgen.NewSine(5.0/2, 0.02, phi, 5.0/2+15.0), // 20mHz
        // refgen.NewSine(5.0/2, 0.03, phi, 5.0/2+15.0), // 30mHz
        // refgen.NewSine(5.0/2, 0.04, phi, 5.0/2+15.0), // 40mHz

        // ===
        // TRI
        // ===
        // 00-20%
        // refgen.NewTriangular(20.0/2, 0.02, phi, 20.0/2), // 20mHz
        // refgen.NewTriangular(20.0/2, 0.03, phi, 20.0/2), // 30mHz
        // refgen.NewTriangular(20.0/2, 0.04, phi, 20.0/2), // 40mHz

        // 00-15%
        // refgen.NewTriangular(15.0/2, 0.02, phi, 15.0/2), // 20mHz
        // refgen.NewTriangular(15.0/2, 0.03, phi, 15.0/2), // 30mHz
        // refgen.NewTriangular(15.0/2, 0.04, phi, 15.0/2), // 40mHz

        // 00-10%
        // refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2), // 20mHz
        // refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2), // 30mHz
        // refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2), // 40mHz

        // 00-05%
        // refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2), // 20mHz
        // refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2), // 30mHz
        // refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2), // 40mHz

        // 05-10%
        // refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+5.0), // 20mHz
        // refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+5.0), // 30mHz
        // refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+5.0), // 40mHz

        // 05-15%
        // refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2+5.0), // 20mHz
        // refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2+5.0), // 30mHz
        // refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2+5.0), // 40mHz

        // 05-20%
        // refgen.NewTriangular(15.0/2, 0.02, phi, 15.0/2+5.0), // 20mHz
        // refgen.NewTriangular(15.0/2, 0.03, phi, 15.0/2+5.0), // 30mHz
        // refgen.NewTriangular(15.0/2, 0.04, phi, 15.0/2+5.0), // 40mHz

        // 10-15%
        // refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+10.0), // 20mHz
        // refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+10.0), // 30mHz
        // refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+10.0), // 40mHz

        // 10-20%
        // refgen.NewTriangular(10.0/2, 0.02, phi, 10.0/2+10.0), // 20mHz
        // refgen.NewTriangular(10.0/2, 0.03, phi, 10.0/2+10.0), // 30mHz
        // refgen.NewTriangular(10.0/2, 0.04, phi, 10.0/2+10.0), // 40mHz

        // 15-20%
        // refgen.NewTriangular(5.0/2, 0.02, phi, 5.0/2+15.0), // 20mHz
        // refgen.NewTriangular(5.0/2, 0.03, phi, 5.0/2+15.0), // 30mHz
        // refgen.NewTriangular(5.0/2, 0.04, phi, 5.0/2+15.0), // 40mHz
    }

    exeInfo := ExeInfo{
        exeTime: 5*60,
        dt: 0.01,
    }

    sensorLength := 140.0 // mm
    loadRef := "1.0" // N
    var err error

    for i, s := range signals {
        // calibration
        err = release(args)
        if err != nil {
            return err
        }

        err = calibrate([]string{loadRef})
        if err != nil {
            return err
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

    return nil
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
    vRef := 5.0
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
    voltageInit := 0.0
    initialVoltageSet := false
    voltageRel := 0.0
    voltageFiltRel := 0.0

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
        voltageRel = (voltage-voltageInit) / voltageInit * 100 // (%)
        voltageFiltRel = (voltageFilt-voltageInit) / voltageInit * 100 // (%)
        
        position, err := hs.Read() // mm
        if err != nil {
            position = prevPositionValue
        }
        if position != -1 {
            prevPositionValue = position
        }
        strain := position / sensorLength * 100 // strain (%)

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
        fmt.Printf("\rTime: %.3f ms / %.3f s | Strain: %.3f", timePerIteration, timeFromStart, strain)
    }}
    fmt.Println("\n\nExperiment finalized")

    return true
}
