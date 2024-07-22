package commands


import (
    "os"
    "os/signal"
    "syscall"
    "time"
    "fmt"

    "goraspio/motor"
    "goraspio/digitalio"
)

var argToDir = map[string]digitalio.PinState {
    "right": digitalio.Low,
    "left": digitalio.High,
}

func move(args []string) error {
    // ARGS
    if len(args) != 1 {
        return fmt.Errorf("[move] wrong number of arguments: expected 1 and got %d", len(args))
    } 
    
    direction, ok := argToDir[args[0]]
    if !ok {
        return fmt.Errorf("[move] wring direction: expected left or right and got %s", args[0])
    }

    // MOTOR
    motor, err := motor.New(13, 1_000, 6)
    if err != nil {
        return fmt.Errorf("[move] error instantiating motor\n%w", err)
    }
    defer motor.Close()

    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    
    ticker := time.NewTicker(time.Millisecond*10)
    defer ticker.Stop()

    fmt.Printf("[move] moving motor to the %s\n       press Ctrl-C to stop\n\n", args[0])

    for {
    select {
    case <- quit:
        fmt.Println("[move] stopping program...")
        return nil

    case <- ticker.C:
        err := motor.WriteRaw(100, direction) 
        if err != nil {
            return fmt.Errorf("[move] error writing to motor\n%w", err)
        }
    }
    }
}
