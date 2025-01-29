package utils

import (
    "os"
    "os/signal"
    "syscall"
    "sync"
)

type KbIntListener struct {
    quit chan os.Signal
    once sync.Once
}

func NewKbIntListener() *KbIntListener {
    stopper := &KbIntListener{
        quit: make(chan os.Signal, 1),
    }
    signal.Notify(stopper.quit, syscall.SIGINT, syscall.SIGTERM)
    return stopper
}

func (gs *KbIntListener) Listen() <-chan os.Signal {
    return gs.quit
}

func (gs *KbIntListener) Stop() {
    gs.once.Do(func() {
        signal.Stop(gs.quit)
        close(gs.quit)
    })
}
