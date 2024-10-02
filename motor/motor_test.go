package motor

import (
	"fmt"
	"math"
	"testing"
	"time"

    "github.com/roboticslab-uc3m/goraspio/digitalio"
    "github.com/roboticslab-uc3m/goraspio/refgen"
)


func TestMotorClosedLoop(t *testing.T) {
    // hall sensor
    hs, err := hallsensor.NewI2C(0x40, 1)
    if err != nil {
        panic(err)
    }
    defer hs.Close()

    // motor
    pwmPinNo := 13
    freq := 2_000
    directionPinNo := 6
    motor, err := New(pwmPinNo, freq, directionPinNo)
    if err != nil {
        t.Fatal(err)
    }
    defer motor.Close()
    
    // refgen
    sine := refgen.NewSine(10/2, 0.04, -math.Pi/2, 10/2)

    rg := refgen.NewRefGen([]refgen.Signal{sine})

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
        
        // reference
        ref := rg.Compute(timeFromStart)

        // position
        position, err := hs.Read()
        if err != nil {
            fmt.Println("Hall sensor failed reading")
        }

        // error
        posError := ref - position
        
        pwmValue, err := motor.Write(posError, dt)
        if err != nil {
            t.Fatal(err)
        }

        timeFromStart = time.Since(programStartTime).Seconds()

        fmt.Printf("\rPosition: %.3f | Reference: %.3f | PWM: %d", position, ref, pwmValue)
    }
}

func TestMotorOpenLoop(t *testing.T) {
    pwmPinNo := 13
    freq := 2_000
    directionPinNo := 6
    motor, err := New(pwmPinNo, freq, directionPinNo)
    if err != nil {
        t.Fatal(err)
    }
    defer motor.Close()

    ref := refgen.NewSine(100, 0.04, -math.Pi/2, 0)

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

        if reference < 0 {
            motor.direction.Write(digitalio.High)
        } else {
            motor.direction.Write(digitalio.Low)
        }

        err := motor.pwm.Write(int(math.Abs(reference)))
        if err != nil {
            t.Fatal(err)
        }

        timeFromStart = time.Since(programStartTime).Seconds()
    }
}

func TestMotorMove(t *testing.T) {
    pwmPinNo := 13
    freq := 10_000
    directionPinNo := 6
    motor, err := New(pwmPinNo, freq, directionPinNo)
    if err != nil {
        t.Fatal(err)
    }
    defer motor.Close()

    // mov := true // right
    mov := false // left

    // ticker
    ticker := time.NewTicker(time.Millisecond * 10)
    defer ticker.Stop()

    // time
    exeTime := 5
    dt := 0.01

    // main loop
    for range int(float64(exeTime)/dt) {
        <-ticker.C

        if mov {
            motor.direction.Write(digitalio.Low)
        } else {
            motor.direction.Write(digitalio.High)
        }

        err := motor.pwm.Write(100)
        if err != nil {
            t.Fatal(err)
        }
    }
}
