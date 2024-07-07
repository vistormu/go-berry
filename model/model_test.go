package model

import (
    "fmt"
    "testing"
)

func TestMain(t *testing.T) {
    contextLength := 1920
    model, err := New(Transformer, contextLength)
    if err != nil {
        t.Fatal(err.Error())
    }
    defer model.Close()

    n_inputs := 5
    inputs := make([][]float32, n_inputs)
    for i := range n_inputs {
        input := make([]float32, contextLength)
        for j := range contextLength {
            input[j] = 0.5*float32(i+1) + float32(j) * 0.01
        }
        inputs[i] = input
    }

    for _, input := range inputs {
        result, err := model.Compute(input)
        if err != nil {
            t.Fatal(err.Error())
        }

        fmt.Println(result)
    }
}
