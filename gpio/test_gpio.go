package gpio

import (
    "testing"
    "time"
    "fmt"
)

func TestDigitalOutput(t *testing.T) {
    pin := Pin(25)
    pin.Output()

    exeTime := 10.0
    dt := 0.01
    ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
    defer ticker.Stop()

    for {
        <- ticker.C
        
        if exeTime <= 0 {
            break
        }

        if exeTime <= 5 {
            pin.High()
        } else {
            pin.Low()
        }
        
        exeTime -= dt
        fmt.Printf("\r%.2f", exeTime)
    }
}
