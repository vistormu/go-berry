package sensor

import (
    "testing"
    "time"
    "fmt"
)

func TestMcp3201(t *testing.T) {
    vs, err := NewMcp3201(5.0, 25, 0.0)
    if err != nil {
        t.Fatal(err)
    }
    defer vs.Close()

    ticker := time.NewTicker(time.Millisecond*10)

    for range 10*100 {
        <- ticker.C

        vr, err := vs.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println(vr.Voltage)
    }
}
