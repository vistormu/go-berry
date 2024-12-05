package gpio

import (
    "testing"
    "time"
    "fmt"
)

func TestDigitalOutput(t *testing.T) {
    defer Close()

    do := NewDigitalOut(25, Low)
    defer do.Close()

    exeTime := 5.0
    dt := 0.01
    ticker := time.NewTicker(time.Duration(dt * float64(time.Second)))
    defer ticker.Stop()

    toggled := false

    for {
        <- ticker.C
        
        if exeTime <= 0 {
            break
        }

        if exeTime <= 2.5 && !toggled {
            do.Toggle()
            toggled = true
        }
        
        exeTime -= dt
        fmt.Printf("\r%.2f, %v", exeTime, do.Read())
    }
}
