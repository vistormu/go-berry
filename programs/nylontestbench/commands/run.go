package commands

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
    "math"

    "gopkg.in/yaml.v3"
    
	"goraspio/client"
	"goraspio/hallsensor"
	"goraspio/motor"
	"goraspio/refgen"
	"goraspio/voltagesensor"
    "goraspio/loadcell"
)

var configPath string = "programs/nylontestbench/config.yaml"

type ExperimentConfig struct {
    Time struct {
        ExeTime int `yaml:"exe_time"`
        Dt float64
    }
    Sensor struct {
        Length float64
        InitialLoad string `yaml:"initial_load"`
    }
    Experiments struct {
        Sinusoidal []struct{
            Freq float64
            MinAmp float64 `yaml:"min_amp"`
            MaxAmp float64  `yaml:"max_amp"`
        }
        Triangular []struct{
            Freq float64
            MinAmp float64 `yaml:"min_amp"`
            MaxAmp float64 `yaml:"max_amp"`
        }
    }
}

func run(args []string) error {
    if len(args) != 0 {
        return fmt.Errorf("[run] wrong number of arguments: expected 0 and got %d", len(args))
    }

    // read config.yaml
    content, err := os.ReadFile(configPath)
    if err != nil {
        return err
    }
    
    ec := ExperimentConfig{}
    err = yaml.Unmarshal(content, &ec)
    if err != nil {
        return err
    }

    // create signals
    phi := -math.Pi/2
    signals := make([]refgen.Signal, len(ec.Experiments.Sinusoidal)+len(ec.Experiments.Triangular))
    for i, s := range ec.Experiments.Sinusoidal {
        signals[i] = refgen.NewSine((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, (s.MaxAmp-s.MinAmp)/2+s.MinAmp)
    }
    for i, s := range ec.Experiments.Triangular {
        signals[i+len(ec.Experiments.Sinusoidal)] = refgen.NewTriangular((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, (s.MaxAmp-s.MinAmp)/2+s.MinAmp)
    }

    for i, s := range signals {
        fmt.Printf("[run] running\n\n")

        // release
        err = release(args)
        if err != nil {
            return err
        }

        // calibration
        err = calibrate([]string{ec.Sensor.InitialLoad})
        if err != nil {
            return err
        }
        
        // experiment
        err := exe([]refgen.Signal{s}, ec.Time.ExeTime, ec.Time.Dt, ec.Sensor.Length)
        if err != nil {
            return err
        }

        // wait for next experiment
        if i != len(signals) - 1 {
            time.Sleep(time.Second*10)
        } 
    }

    return nil
}

func exe(signals []refgen.Signal, exeTime int, dt float64, sensorLength float64) error {
    // ==========
    // COMPONENTS
    // ==========
    // Motor
    pwmPinNo := 13
    freq := 4_500
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        return err
    }
    defer motor.Close()

    // Voltage Sensor
    vRef := 5.0
    voltageSensorhipSelectNo := 25
    vs, err := voltagesensor.New(vRef, voltageSensorhipSelectNo)
    if err != nil {
        return err
    }
    defer vs.Close()

    // Hall Sensor
    hs, err := hallsensor.NewI2C(0x40, 1)
    if err != nil {
        return err
    }
    defer hs.Close()

    // Load cell
    lc, err := loadcell.New(24)
    if err != nil {
        return err
    }
    defer lc.Close()

    // Client
    ip := "10.118.90.193"
    port := 8080
    c, err := client.New(ip, port)
    if err != nil {
        return err
    }
    defer c.Close()

    // Reference generator
    rg := refgen.NewRefGen(signals)

    // ========
    // CHANNELS
    // ========
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    ticker := time.NewTicker(time.Duration(dt*float64(time.Second)))
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

    for range int(float64(exeTime)/dt) {
    select {
    case <- quit:
        return fmt.Errorf("[run] keyboard interrupt")
    case <-ticker.C:
        loopStartTime := time.Now()

        // READ
        voltage, voltageFilt, err := vs.Read() // V
        if err != nil {
            return err
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
            return err
        }

        // REFERENCE
        ref := rg.Compute(time.Since(programStartTime).Seconds()) // strain

        // ACTUATE
        _, err = motor.Write(ref-strain, dt)
        if err != nil {
            return err
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
            return err
        }

        // TIME
        timePerIteration := time.Since(loopStartTime).Seconds()*1000
        timeFromStart = time.Since(programStartTime).Seconds()

        // PRINT
        fmt.Printf("\rtime per iteration: %.3f ms | execution time: %.0f s", timePerIteration, timeFromStart)
    }}
    fmt.Println("\n\nExperiment finalized")

    return nil
}
