package commands


import (
    "fmt"
    "os"
)

type commandFunc func(args []string) error
var commandStrToFunc = map[string]commandFunc {
    "move": move,
    "calibrate": calibrate,
    "release": release,
    "run": run,
}

func exit(err error) {
    fmt.Println(err.Error())
    os.Exit(1)
}

func Execute(args []string) {
    if len(args) == 0 {
        exit(fmt.Errorf("[testbench] wrong number of args: expected more than one"))
    }

    command, ok := commandStrToFunc[args[0]]
    if !ok {
        exit(fmt.Errorf("[testbench] command not found: %s", args[0]))
    }

    err := command(args[1:])
    if err != nil {
        exit(err)
    }
}
