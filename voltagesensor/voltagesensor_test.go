package voltagesensor

import (
    "testing"
    "time"
    "fmt"
)

func TestVoltageSensor(t *testing.T) {
    vs, err := New(5.0, 25)
    if err != nil {
        t.Fatal(err)
    }
    defer vs.Close()

    ticker := time.NewTicker(time.Millisecond*10)

    for range 10*100 {
        <- ticker.C

        voltage, _, err := vs.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println(voltage)
    }
}
