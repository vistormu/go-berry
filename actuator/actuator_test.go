package actuator

import (
    "testing"
    "time"
    "fmt"
)

func Test17hs4401(t *testing.T) {
    m, err := NewStepMotor17hs4401(18, 500, 15)
    if err != nil {
        t.Fatal(err)
    }
    defer m.Close()

    ticker := time.NewTicker(time.Millisecond*10)
    defer ticker.Stop()

    startTime := time.Now()

    for range 1_000 {
        <- ticker.C

        fmt.Println(time.Since(startTime).Seconds())

        pwm := -100.0
        err = m.Write(pwm)
        if err != nil {
            t.Fatal(err)
        }
    }
}
