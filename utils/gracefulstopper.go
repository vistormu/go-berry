package utils

import (
    "os"
    "os/signal"
    "syscall"
)

type GracefulStopper struct {
    quit chan os.Signal
}

func NewGracefulStopper() *GracefulStopper {
    stopper := &GracefulStopper{
        quit: make(chan os.Signal, 1),
    }
    signal.Notify(stopper.quit, syscall.SIGINT, syscall.SIGTERM)
    return stopper
}

func (gs *GracefulStopper) ListenForShutdown() <-chan os.Signal {
    return gs.quit
}
