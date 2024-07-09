package hallsensor

import (
    "fmt"
    "time"
    "testing"
    "goraspio/digitalio"
)

func TestHallSensorSpi(t *testing.T) {
    hs, err := NewSpi(24)
    if err != nil {
        t.Fatal(err)
    }
    defer hs.Close()

    m, err := digitalio.NewPwm(13, 2_000)
    if err != nil {
        t.Fatal(err)
    }
    defer m.Close()

    direction := digitalio.NewDigitalOut(6, digitalio.Low)
    defer direction.Close()
    
    // slice from 0 to 100 and back
    length := 200
    values := make([]int, length)
    for i := range length {
        if i < length/2 {
            values[i] = i
        } else {
            values[i] = length - i
        }
    }

    ticker := time.NewTicker(time.Millisecond*10)

    for range 20 {
        direction.Toggle()
        // direction.Write(digitalio.Low)
        // direction.Write(digitalio.High)
        for _, v := range values {
            <-ticker.C

            err := m.Write(v)
            if err != nil {
                t.Fatal("couldnt write")
            }

            position, err := hs.Read()
            if err != nil {
                t.Fatal(err)
            }

            fmt.Printf("%f\n", position)
        }
    }
}

func TestHallSensorI2C(t *testing.T) {
    var address byte = 0x40
    line := 1
    hs, err := NewI2C(address, line)
    if err != nil {
        t.Fatal(err)
    }
    defer hs.Close()

    ticker := time.NewTicker(time.Millisecond * 10)

    ti := 15*1000/10

    for range ti {
        <-ticker.C
        position, err := hs.Read()
        if err != nil {
            t.Fatal(err)
        }

        fmt.Printf("\n\nPosition: %d\n\n", position)
    }
}
