package main

import (
    "os"
    "goraspio/programs/nylontestbench/commands"
)

func main() {
    commands.Execute(os.Args[1:])
}
