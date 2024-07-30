package commands


import (
    "fmt"
    "os"
)

func exit(err error) {
    fmt.Println(err.Error())
    os.Exit(1)
}

func Execute(args []string) {
    if len(args) == 0 {
        exit(fmt.Errorf("[testbench] wrong number of args: expected more than one"))
    }

    var err error
    switch args[0] {
    case "move":
        err = move(args[1:])
    case "calibrate":
        _, err = calibrate(args[1:])
    case "release":
        err = release(args[1:])
    case "run":
        err = run(args[1:])
    case "reach":
        err = reach(args[1:])
    default:
        exit(fmt.Errorf("[testbench] command not found: %s", args[0]))
    }

    if err != nil {
        exit(err)
    }
}
