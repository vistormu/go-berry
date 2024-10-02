package utils

import (
    "fmt"
    "testing"
    "time"
)

func TestWindow(t *testing.T) {
    w := NewWindow(10)
    fmt.Println(w.Data)
    
    w.Append(1.0)

    fmt.Println(w.Data)
}

func TestPrint(t *testing.T) {
    for i := 0; i < 100; i++ {
        fmt.Printf("\rProcessing %d%% complete", i)
        time.Sleep(100 * time.Millisecond) // Simulate work
    }
    fmt.Println("\rProcessing 100% complete")
}

func TestGracefulStopper(t *testing.T) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    stopper := NewGracefulStopper()

    done := false
    for !done {
    select{
    case <-ticker.C:
        fmt.Println("Tick")

    case <-stopper.ListenForShutdown():
        fmt.Println("Shutting down")
        done = true
        break
    }
    }

    fmt.Println("Done")
}
