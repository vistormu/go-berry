package commands

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopkg.in/yaml.v3"

	"goraspio/client"
	"goraspio/hallsensor"
	"goraspio/loadcell"
	"goraspio/model"
	"goraspio/motor"
	"goraspio/refgen"
	"goraspio/utils"
	"goraspio/voltagesensor"
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
    // signals := make([]refgen.Signal, len(ec.Experiments.Sinusoidal)+len(ec.Experiments.Triangular))
    // for i, s := range ec.Experiments.Sinusoidal {
    //     signals[i] = refgen.NewSine((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, (s.MaxAmp-s.MinAmp)/2+s.MinAmp)
    // }
    // for i, s := range ec.Experiments.Triangular {
    //     index := i+len(ec.Experiments.Sinusoidal)
    //     signals[index] = refgen.NewTriangular((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, s.MinAmp)
    // }

    // for _, s := range ec.Experiments.Sinusoidal {
    //     fmt.Printf("[run] running sinusoidal with min amp: %.1f, max amp: %.1f, freq: %.2f\n\n", s.MinAmp, s.MaxAmp, s.Freq)

    //     signal := refgen.NewSine((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, (s.MaxAmp-s.MinAmp)/2+s.MinAmp)

    //     // release
    //     err = release(args)
    //     if err != nil {
    //         return err
    //     }

    //     // calibration
    //     voltageInit, err := calibrate([]string{ec.Sensor.InitialLoad})
    //     if err != nil {
    //         return err
    //     }

    //     // reach
    //     err = reach([]string{fmt.Sprintf("%.2f", s.MinAmp), fmt.Sprintf("%.2f", ec.Sensor.Length)})
    //     if err != nil {
    //         return err
    //     }
        
    //     // experiment
    //     err = exe([]refgen.Signal{signal}, ec.Time.ExeTime, ec.Time.Dt, ec.Sensor.Length, voltageInit, s.MinAmp)
    //     if err != nil {
    //         return err
    //     }

    //     // wait for next experiment
    //     time.Sleep(time.Second*10)
    // }

    // for _, s := range ec.Experiments.Triangular {
    //     fmt.Printf("[run] running triangular with min amp: %.1f, max amp: %.1f, freq: %.2f\n\n", s.MinAmp, s.MaxAmp, s.Freq)

    //     signal := refgen.NewTriangular((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, s.MinAmp)

    //     // release
    //     err = release(args)
    //     if err != nil {
    //         return err
    //     }

    //     // calibration
    //     voltageInit, err := calibrate([]string{ec.Sensor.InitialLoad})
    //     if err != nil {
    //         return err
    //     }

    //     // reach
    //     err = reach([]string{fmt.Sprintf("%.2f", s.MinAmp), fmt.Sprintf("%.2f", ec.Sensor.Length)})
    //     if err != nil {
    //         return err
    //     }
        
    //     // experiment
    //     err = exe([]refgen.Signal{signal}, ec.Time.ExeTime, ec.Time.Dt, ec.Sensor.Length, voltageInit, s.MinAmp)
    //     if err != nil {
    //         return err
    //     }

    //     // wait for next experiment
    //     time.Sleep(time.Second*10)
    // }

    //     signals[i] = refgen.NewSine((s.MaxAmp-s.MinAmp)/2, s.Freq, phi, (s.MaxAmp-s.MinAmp)/2+s.MinAmp)
    mixed_10_20 := []refgen.Signal{
        refgen.NewSine(2.5/2, 0.01, phi, 2.5/2),
        refgen.NewSine(2.5/2, 0.02, phi, 2.5/2),
    }
    mixed_15_25 := []refgen.Signal{
        refgen.NewSine(2.5/2, 0.015, phi, 2.5/2),
        refgen.NewSine(2.5/2, 0.025, phi, 2.5/2),
    }
    mixed_20_30 := []refgen.Signal{
        refgen.NewSine(2.5/2, 0.02, phi, 2.5/2),
        refgen.NewSine(2.5/2, 0.03, phi, 2.5/2),
    }

    signalsExp := [][]refgen.Signal{
        mixed_10_20,
        mixed_15_25,
        mixed_20_30,
    }

    for i, s := range signalsExp {
        fmt.Printf("[run] running experiment %d", i)

        // release
        err = release(args)
        if err != nil {
            return err
        }

        // calibration
        voltageInit, err := calibrate([]string{ec.Sensor.InitialLoad})
        if err != nil {
            return err
        }

        // reach
        err = reach([]string{fmt.Sprintf("%.2f", 0.0), fmt.Sprintf("%.2f", ec.Sensor.Length)})
        if err != nil {
            return err
        }
        
        // experiment
        err = exe(s, ec.Time.ExeTime, ec.Time.Dt, ec.Sensor.Length, voltageInit, 0.0)
        if err != nil {
            return err
        }

        // wait for next experiment
        time.Sleep(time.Second*10)
    }


    err = release(args)
    if err != nil {
        return err
    }

    return nil
}

func exe(signals []refgen.Signal, exeTime int, dt float64, sensorLength float64, voltageInit float64, strainInit float64) error {
    // ==========
    // COMPONENTS
    // ==========
    // Motor
    pwmPinNo := 13
    freq := 2000
    dirPinNo := 6
    motor, err := motor.New(pwmPinNo, freq, dirPinNo)
    if err != nil {
        return err
    }
    defer motor.Close()

    // Voltage Sensor
    vRef := 5.0
    voltageSensorhipSelectNo := 25
    vs, err := voltagesensor.New(vRef, voltageSensorhipSelectNo, voltageInit)
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

    // Models
    strainModel, err := model.New("TransformerRegressorStrain", 1920)
    if err != nil {
        return err
    }
    defer strainModel.Close()
    voltageWindow := utils.NewWindow(1920)
    strainModelKf := utils.NewKalmanFilter(0.05, 20, 0.01)
    strainModelKf.SetInitialEstimate(0)

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
    data := make(map[string]float64)

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
        vr, err := vs.Read() // V
        if err != nil {
            return err
        }
        
        position, _ := hs.Read() // mm
        strain := position / sensorLength * 100 + strainInit // strain (%)

        load, loadFilt, err := lc.Read() // N
        if err != nil {
            return err
        }

        // PREDICT
        voltageWindow.Append(vr.VoltageFiltRel/100)
        strainPred, err := strainModel.Compute(voltageWindow.Data)
        if err != nil {
            return err
        }
        strainPred = strainModelKf.Compute(strainPred)

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

        data["voltage"] = vr.Voltage
        data["voltage_filt"] = vr.VoltageFilt
        data["voltage_rel"] = vr.VoltageRel
        data["voltage_filt_rel"] = vr.VoltageFiltRel

        data["reference"] = ref
        data["position"] = position

        data["strain"] = strain
        data["strain_pred"] = strainPred

        err = c.Send(data)
        if err != nil {
            return err
        }

        // TIME
        timePerIteration := time.Since(loopStartTime).Seconds()*1000
        timeFromStart = time.Since(programStartTime).Seconds()

        // PRINT
        fmt.Printf("\rtime per iteration: %.3f ms | execution time: %.0f/%d s", timePerIteration, timeFromStart, exeTime)
    }}
    fmt.Println("\n\nExperiment finalized")

    return nil
}
