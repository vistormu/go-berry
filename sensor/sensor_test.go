package sensor

import (
    "testing"
    "time"
    "fmt"
)

func TestMcp3201(t *testing.T) {
    vs, err := NewMcp3201(5.0, 25)
    if err != nil {
        t.Fatal(err)
    }
    defer vs.Close()

    ticker := time.NewTicker(time.Millisecond*10)

    for range 10*100 {
        <- ticker.C

        v, err := vs.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println(v)
    }
}

func TestAs5048b(t *testing.T) {
    hs, err := NewAs5048b(0x40, 1)
    if err != nil {
        t.Fatal(err)
    }
    defer hs.Close()
    
    exeTime := 30.0
    dt := 0.1
    ticker := time.NewTicker(time.Duration(dt*float64(time.Second)))

    for range int(exeTime/dt) {
        <- ticker.C

        angle, err := hs.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Println(angle)
    }
}
