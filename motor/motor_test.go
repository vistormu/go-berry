package motor

import (
    "time"
    "math"
    "testing"
    "goraspio/refgen"
)


func TestMotor(t *testing.T) {
    pwmPinNo := 13
    freq := 2_000
    directionPinNo := 6
    motor, err := New(pwmPinNo, freq, directionPinNo)
    if err != nil {
        t.Fatal(err)
    }
    defer motor.Close()

    ref := refgen.NewSine(10/2, 0.04, -math.Pi/2, 10/2)

    // ticker
    ticker := time.NewTicker(time.Millisecond * 10)
    defer ticker.Stop()

    // time
    exeTime := 25
    dt := 0.01
    programStartTime := time.Now()
    timeFromStart := 0.0

    // main loop
    for range int(float64(exeTime)/dt) {
        <-ticker.C

        reference := ref.Compute(timeFromStart)
        
        err := motor.Write(reference, reference-0.1)
        if err != nil {
            t.Fatal("couldnt write")
        }

        timeFromStart = time.Since(programStartTime).Seconds()
    }
}
