package errors

import (
    "fmt"
    "strings"
    "reflect"
)

const (
    END = "\x1b[0m"
    ITEM = "\n   |> "
    FULL = ITEM + "full error:\n\n%v"
)

type ErrorType interface {
    String() string
}

// ====
// GPIO
// ====
type GpioError string
const (
    GPIO_INIT GpioError = "error initializing gpio pins" + END + ITEM + "registers: %v" + FULL
    GPIO_CLOSE GpioError = "error closing gpio" + END + FULL
    GPIO_BASE GpioError = "error reading base address" + END + FULL
)
func (e GpioError) String() string {
    return string(e)
}

// ===
// PWM
// ===
type PwmError string
const (
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
    I2C_OPEN I2CError = "error opening i2c channel" + END + FULL
    I2C_READ I2CError = "error reading from i2c" + END + ITEM + "register: %v" + FULL
    I2C_WRITE I2CError = "error writing to i2c" + END + ITEM + "register: %v" + FULL
    I2C_CLOSE I2CError = "error closing i2c" + END + FULL
)
func (e I2CError) String() string {
    return string(e)
}

// ======
// client
// ======
type ClientError string
const (
    CONNECTION ClientError = "could not connect to the server" + END + FULL
    CLIENT_SEND ClientError = "could not send data" + END + FULL
    CLIENT_JSON ClientError = "could not encode data" + END + FULL
    CLIENT_CLOSE ClientError = "could not close connection" + END + FULL
)
func (e ClientError) String() string {
    return string(e)
}

var stageMessages = map[reflect.Type]string{
    reflect.TypeOf(GpioError("")): "|gpio error| ",
    reflect.TypeOf(PwmError("")): "|pwm error| ",
    reflect.TypeOf(SpiError("")): "|spi error| ",
    reflect.TypeOf(I2CError("")): "|i2c error| ",
    reflect.TypeOf(ClientError("")): "|client error| ",
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

func Must[T any](obj T, err error) T {
    if err != nil {
        panic(err)
    }

    return obj
}
