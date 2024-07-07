package main

import (
    "os"
    "os/signal"
    "syscall"
    "time"
    "fmt"
    "math"
    "goraspio/digitalio"
    "goraspio/voltagesensor"
    "goraspio/hallsensor"
    "goraspio/client"
    "goraspio/controller"
    "goraspio/refgen"
    "goraspio/model"
    "goraspio/utils"
)

type ModelInfo struct {
    modelType model.ModelType
    contextLength int
}

type ExeInfo struct {
    exeTime int
    dt float64
}

func main() {
    amp := 20.0/2
    freq := 0.03
    phi := -math.Pi/2
    offset := amp
    sine := refgen.NewSine(amp, freq, phi, offset)

    signals := make([]refgen.Signal, 1)
    signals[0] = sine

    modelInfo := ModelInfo{
        modelType: model.Transformer,
        contextLength: 1920,
    }
    exeInfo := ExeInfo{
        exeTime: 67,
        dt: 0.01,
    }

    exe(signals, modelInfo, exeInfo)
}

func exe(signals []refgen.Signal, modelInfo ModelInfo, exeInfo ExeInfo) {
    // ==========
    // COMPONENTS
    // ==========
    // PWM
    pinNo := 18
    freq := 500
    master, err := digitalio.NewPwm(pinNo, freq)
    if err != nil {
        panic(err)
    }
    defer master.Close()
    fmt.Println("master pwm connected successfully")

    slave, err := digitalio.NewPwm(13, freq)
    if err != nil {
        panic(err)
    }
    defer slave.Close()
    fmt.Println("slave pwm connected successfully")

    // Voltage Sensor
    var vRef float32 = 5.0
    chipSelectNo := 6
    vs, err := voltagesensor.New(vRef, chipSelectNo)
    if err != nil {
        panic(err)
    }
    defer vs.Close()
    fmt.Println("voltage sensor connected successfully")

    // Hall Sensor
    var address byte = 0x40
    line := 1
    hs, err := hallsensor.NewI2C(address, line)
    if err != nil {
        panic(err)
    }
    defer hs.Close()
    fmt.Println("hall sensor connected successfully")

    // Client
    ip := "163.117.150.172"
    port := 8080
    c, err := client.New(ip, port)
    if err != nil {
        panic(err)
    }
    defer c.Close()
    fmt.Println("client connected successfully")

    // Controller
    var kp  float32 = 0.02
    var ki  float32 = 0.01
    var kd  float32 = 0.01
    var alpha float32 = 0.9
    bounds := [2]float32{-1000.0, 1000.0}
    masterPID := controller.NewPID(kp, ki, kd, alpha, bounds)
    slavePID := controller.NewPID(kp/2, ki/2, kd/2, alpha, bounds)

    // Reference generator
    rg := refgen.NewRefGen(signals)
    if err != nil {
        panic(err)
    }

    // model
    m, err := model.New(modelInfo.modelType, modelInfo.contextLength)
    if err != nil {
        panic(err)
    }
    defer m.Close()
    fmt.Println("model initialized")

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
    modelInput := utils.NewWindow(modelInfo.contextLength)
    
    var prevRef float32 = 0.0

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
            panic(err)
        }

        // MODEL
        modelInput.Append(filteredVoltage)
        modelOutput, err := m.Compute(modelInput.Data)
        if err != nil {
            panic(err)
        }

        // REFERENCE
        ref := rg.Compute(timeFromStart)

        // CONTROL
        error := ref - float32(position)
        
        masterControl := 0
        slaveControl := 0
        if error < -1000 && ref - prevRef < 0 {
            slaveControl = int(slavePID.Compute(-error, float32(exeInfo.dt)))
        } else {
            masterControl = int(masterPID.Compute(error, float32(exeInfo.dt)))
        }

        // ACTUATE
        err = master.Write(masterControl)
        if err != nil {
            panic(err)
        }
        err = slave.Write(slaveControl)
        if err != nil {
            panic(err)
        }

        // SEND
        data["time"] = timeFromStart

        data["master_control"] = masterControl
        data["control"] = slaveControl

        data["voltage"] = voltage
        data["filtered_voltage"] = filteredVoltage

        data["position"] = position
        data["reference"] = ref
        data["model"] = modelOutput

        err = c.Send(data)
        if err != nil {
            panic(err)
        }

        // TIME
        timePerIteration := time.Since(loopStartTime).Seconds()*1000
        timeFromStart = time.Since(programStartTime).Seconds()

        prevRef = ref

        // PRINT
        fmt.Printf("\rTime: %.3f ms / %.3f s | Voltage: %.3f | Position: %d | Control: %d", timePerIteration, timeFromStart, voltage, position, int(masterControl))
    }}
    fmt.Println("\nProgram finalized")
}
