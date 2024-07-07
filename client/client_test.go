package client

import (
    "testing"
    "time"
)

func TestMain(t *testing.T) {
    ip := "163.117.150.172"
    port := 8080
    c, err := New(ip, port)
    if err != nil {
        t.Fatal("couldnt open")
    }
    defer c.Close()

    ticker := time.NewTicker(time.Millisecond)
    defer ticker.Stop()

    data := map[string]any{
        "resistance": 100,
        "position": 1_000,
        "reference": 20_000,
        "control": 1.0,
        "time": 0.0,
    }

    for i := range 20_000 {
        <- ticker.C
        data["time"] = 0.001*float64(i)
        data["position"] = data["position"].(int) + 1
        c.Send(data)
    }
}
