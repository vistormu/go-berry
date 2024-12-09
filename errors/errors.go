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

// ===
// I2C
// ===
type I2CError string
const (
    // pwm errors
    I2C_OPEN I2CError = "error opening i2c channel" + END + ITEM + "full error:\n\n%v"
    I2C_READ I2CError = "error reading from i2c" + END + ITEM + "register: %v" + ITEM + "full error:\n\n%v"
    I2C_WRITE I2CError = "error writing to i2c" + END + ITEM + "register: %v" + ITEM + "full error:\n\n%v"
)
func (e I2CError) String() string {
    return string(e)
}

var stageMessages = map[reflect.Type]string{
    reflect.TypeOf(PwmError("")): "|pwm error| ",
    reflect.TypeOf(SpiError("")): "|spi error| ",
    reflect.TypeOf(I2CError("")): "|i2c error| ",
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
