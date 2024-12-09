package errors

import (
    "fmt"
    "strings"
    "reflect"
)

const (
    END = "\x1b[0m"
    ITEM = "\n   |> "
)

type ErrorType interface {
    String() string
}

// ===
// PWM
// ===
type PwmError string
const (
    // pwm errors
    PWM_PIN PwmError = "wrong pwm pin" + END + ITEM + "got: %v" + ITEM + "available pwm pins: 12, 13, 40, 41, 45 | 18, 19"
)
func (e PwmError) String() string {
    return string(e)
}

// ===
// SPI
// ===
type SpiError string
const (
    // pwm errors
    SPI_ROOT SpiError = "error mapping spi registers" + END + ITEM + "are you root?"
)
func (e SpiError) String() string {
    return string(e)
}

var stageMessages = map[reflect.Type]string{
    reflect.TypeOf(PwmError("")): "|pwm error| ",
    reflect.TypeOf(SpiError("")): "|spi error| ",
}

type Error struct {
    message string
}

func New(errorType ErrorType, args ...any) error {
    stageMessage := stageMessages[reflect.TypeOf(errorType)]
    errorMessage := errorType.String()

    message := "\x1b[31m-> " + stageMessage + errorMessage + "\n"
    n := strings.Count(message, "%v")

    if len(args) != n {
        panic(fmt.Sprintf("expected %v arguments, got %v", n, len(args)))
    }

    message = fmt.Sprintf(message, args...)

    return Error{message}
}

func (e Error) Error() string {
    return e.message
}
