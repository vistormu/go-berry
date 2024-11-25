package utils

import (
    "os"
    "os/signal"
    "syscall"
    "sync"
)

type GracefulStopper struct {
    quit chan os.Signal
    once sync.Once
}

func NewGracefulStopper() *GracefulStopper {
    stopper := &GracefulStopper{
        quit: make(chan os.Signal, 1),
    }
    signal.Notify(stopper.quit, syscall.SIGINT, syscall.SIGTERM)
    return stopper
}

func (gs *GracefulStopper) Listen() <-chan os.Signal {
    return gs.quit
}

func (gs *GracefulStopper) Stop() {
    gs.once.Do(func() {
        signal.Stop(gs.quit)
        close(gs.quit)
    })
}
