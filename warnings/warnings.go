package warnings

import (
    "fmt"
    "strings"
    "reflect"
)

const (
    END = "\x1b[0m"
    ITEM = "\n   |> "
)

var enabled = true

type WarningType interface {
    String() string
}

type PwmWarning string
const (
    // pwm errors
    DUTY_CYCLE PwmWarning = "duty cycle must be in the range 0-100" + END + ITEM + "got: %v" + ITEM + "value clipped to: %v"
    FREQUENCY PwmWarning = "frequency cannot be negative or greater than the sampling rate" + END + ITEM + "got: %v" + ITEM + "value clipped to: %v"
)

func (e PwmWarning) String() string {
    return string(e)
}

var stageMessages = map[reflect.Type]string{
    reflect.TypeOf(PwmWarning("")): "|pwm warning| ",
}

func New(warningType WarningType, args ...any) {
    stageMessage := stageMessages[reflect.TypeOf(warningType)]
    errorMessage := warningType.String()

    message := "\x1b[33m-> " + stageMessage + errorMessage + "\n"
    n := strings.Count(message, "%v")

    if len(args) != n {
        panic(fmt.Sprintf("expected %v arguments, got %v", n, len(args)))
    }

    if enabled {
        fmt.Printf(message, args...)
    }
}

func Disable() {
    enabled = false
}
