package refgen

import (
	"goraspio/client"
	"testing"
	"time"
    "math"
)

func TestMain(t *testing.T) {
    // references
    signals := make([][]Signal, 0)

    amp := 20.0/2

    // 1: sine
    sine := NewSine(amp, 0.03, -math.Pi/2, amp)
    signals = append(signals, []Signal{sine})

    // 2: triangular
    tri := NewTriangular(amp, 0.03, -math.Pi/2, amp)
    signals = append(signals, []Signal{tri})

    // 3: square
    sqr := NewSquare(amp, 0.03, 0.0, amp)
    signals = append(signals, []Signal{sqr})

    // 4: mixed sins
    sin1 := NewSine(amp/2, 0.02, -math.Pi/2, amp/2)
    sin2 := NewSine(amp/2, 0.04, -math.Pi/2, amp/2)
    signals = append(signals, []Signal{sin1, sin2})

    // client
    c, err := client.New("10.118.90.193", 8080)
    if err != nil {
        t.Fatal(err.Error())
    }
    defer c.Close()

    // times
    dt := time.Millisecond*10

    ticker := time.NewTicker(dt)
    defer ticker.Stop()

    for _, s := range signals {
        time.Sleep(time.Second*5)

        // data
        data := map[string]float64 {
            "time": 0.0,

            "master_control": 0.0,
            "control": 0.0,

            "resistance": 0.0,

            "position": 0.0,
            "reference": 0.0,
            "model": 0.0,
        }

        rg := NewRefGen(s)

        for range int(40/dt.Seconds()) {
            <- ticker.C

            data["reference"] = float64(rg.Compute(data["time"]))

            err := c.Send(data)
            if err != nil {
                t.Fatal(err.Error())
            }
            
            data["time"] += dt.Seconds()
        }
    }
}
