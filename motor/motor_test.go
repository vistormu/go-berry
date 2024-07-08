package motor

import (
    "time"
    "testing"
)

func generateSlice(size int) []float64 {
    if size <= 1 {
        return []float64{-1} // handle edge case where size is 1 or less
    }

    slice := make([]float64, size)
    step := 2.0 / float64(size-1) // calculate the step size
    for i := 0; i < size; i++ {
        slice[i] = -1 + step*float64(i)
    }
    return slice
}

func TestMotor(t *testing.T) {
    pwmPinNo := 13
    freq := 2_000
    directionPinNo := 6
    motor, err := New(pwmPinNo, freq, directionPinNo)
    if err != nil {
        t.Fatal(err)
    }
    defer motor.Close()

    // slice from -1.0 to 1.0
    values := generateSlice(200)

    // ticker
    ticker := time.NewTicker(time.Millisecond * 10)
    defer ticker.Stop()

    // main loop
    for range 5 {
        for _, v := range values {
            <-ticker.C
            
            err := motor.Write(v)
            if err != nil {
                t.Fatal("couldnt write")
            }
        }
    }
}
