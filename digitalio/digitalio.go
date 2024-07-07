package digitalio

import (
    "github.com/stianeikeland/go-rpio/v4"
)

func init() {
    err := rpio.Open()
    if err != nil {
        panic(err)
    }
}

func Close() {
    rpio.Close()
}
