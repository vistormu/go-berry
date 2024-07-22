package loadcell

import (
    "testing"
    "time"
    "fmt"
)

func TestVoltageSensor(t *testing.T) {
    lc, err := New(24)
    if err != nil {
        t.Fatal(err)
    }
    defer lc.Close()

    ticker := time.NewTicker(time.Millisecond*10)

    for range 10*100 {
        <- ticker.C

        load, _, err := lc.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println(load)
    }
}
