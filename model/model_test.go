package model

import (
    "fmt"
    "testing"
    "time"
    "github.com/vistormu/goraspio/num"
)

func TestMain(t *testing.T) {
    contextLength := 1920
    model, err := New("TransformerRegressorStrain", contextLength)
    if err != nil {
        t.Fatal(err.Error())
    }
    defer model.Close()

    n_inputs := 10_000
    inputs := make([][]float64, n_inputs)
    for i := range n_inputs {
        input := make([]float64, contextLength)
        for j := range contextLength {
            input[j] = 0.5*float64(i+1) + float64(j) * 0.01
        }
        inputs[i] = input
    }

    times := make([]float64, n_inputs)
    for i, input := range inputs {
        before := time.Now()
        _, err := model.Compute(input)
        times[i] = time.Since(before).Seconds()*1000
        if err != nil {
            t.Fatal(err.Error())
        }

        // fmt.Println(result)
    }
    fmt.Printf("time: %.2f +/- %.2f", num.Mean(times), num.StdDev(times))
}
