package actuator

import (
    "fmt"
    "testing"
    "time"

    "github.com/vistormu/goraspio/utils"
)

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

    stopper := utils.NewGracefulStopper()

    speed := 70.0
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

            err := m.Write(speed)
            if err != nil {
                t.Fatalf(err.Error())
            }
            
            exeTime -= dt
            fmt.Printf("\r%.2f", exeTime)
        }
    }
}
