package sensor

import (
    "testing"
    "time"
    "fmt"
)


func TestSensor(t *testing.T) {
    sensorName := "mcp3201"
    var sensor Sensor
    var err error
    switch sensorName {
        case "mcp3201":
            sensor, err = NewMcp3201(5.0, 24)
        // case "ems20":
        //     sensor, err = NewEms20(23)
        // case "as5048a":
        //     sensor, err = NewAs5048a(25)
        // case "as5311":
        //     sensor, err = NewAs5311(25)
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

        value := sensor.Read()

        fmt.Println(value)
    }
}
