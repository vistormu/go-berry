package digitalio


import (
    "fmt"
    "time"
    "testing"
)

func TestPWM(t *testing.T) {
    pwm, err := NewPwm(13, 2000)
    if err != nil {
        t.Fatal(err)
    }
    defer pwm.Close()

    direction := NewDigitalOut(6, Low)
    defer direction.Close()
    
    // slice from 0 to 100 and back
    values := make([]int, 200)
    for i := range 200 {
        if i < 100 {
            values[i] = i
        } else {
            values[i] = 200 - i
        }
    }

    // ticker
    ticker := time.NewTicker(time.Millisecond * 10)

    // main loop
    for i := range 5 {
        fmt.Println(i)
        for _, v := range values {
            <-ticker.C
            
            if i % 2 == 0 {
                direction.Write(High)
            } else {
                direction.Write(Low)
            }

            err := pwm.Write(v)
            if err != nil {
                t.Fatal("couldnt write")
            }
        }
    }
}

