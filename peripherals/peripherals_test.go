package peripherals

import (
    "testing"
    "time"
    "fmt"

    "github.com/vistormu/go-berry/utils"
)


func TestSensor(t *testing.T) {
    // sensorName := "mcp3201"
    sensorName := "nse5310"
    var sensor Sensor
    var err error
    switch sensorName {
    case "mcp3201":
        sensor, err = NewMcp3201(5.0, 24)
    case "nse5310":
        sensor, err = NewNse5310(0x40, 1)
    default:
        t.Fatal("unknown sensor")
    }
    if err != nil {
        t.Fatal(err)
    }
    defer sensor.Close()

    exeTime := 10.0
    dt := 0.001
    ticker := time.NewTicker(time.Duration(dt*float64(time.Second)))
    defer ticker.Stop()

    for range int(exeTime/dt) {
        <- ticker.C

        startTime := time.Now()

        value, _ := sensor.Read()

        finish := time.Since(startTime).Seconds() * 1000

        fmt.Printf("\r%.2f V | %.2f ms", value, finish)
    }
}

func Test17hs4401Displacement(t *testing.T) {
    m, err := NewStepMotor17hs4401(18, 15, 10, 10_000)
    if err != nil {
        t.Fatal(err)
    }
    defer m.Close()

    exeTime := 10.0
    dt := 0.01
    ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
    defer ticker.Stop()

    stopper := utils.NewKbIntListener()
    defer stopper.Stop()

    speed := 20.0
    reversed := false
    pulses := int(exeTime / dt)

    for i := range pulses {
        select {
        case <-stopper.Listen():
            return

        case <-ticker.C:
            if i >= pulses / 2 && !reversed {
                speed = -speed
                reversed = true
            }

            m.Write(speed)
            
            exeTime -= dt
            fmt.Printf("\r%.2f", exeTime)
        }
    }
}
