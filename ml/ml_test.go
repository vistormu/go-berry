package ml

import (
    "fmt"
    "testing"
    "time"
)

func TestLocal(t *testing.T) {
    nInputs := 4
    nOutputs := 3
    model, err := NewLocalModel("fnn.onnx", nInputs, nOutputs)
    if err != nil {
        t.Fatal(err.Error())
    }
    defer model.Close()

    input := []float64{0.2, 0.1, 0.12, 0.0}
    output, err := model.Compute(input)
    if err != nil {
        t.Fatal(err)
    }
    fmt.Println(output)
}

func TestRemote(t *testing.T) {
    model, _ := NewRemoteModel(
        "145.94.123.92",
        8080,
        "predict",
        "output",
        []string{"s0", "s1", "s2", "s3"},
    )

    start := time.Now()
    modelOutput, err := model.Compute([]float64{0.2, 0.1, 0.12, 0.0})
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(time.Since(start))
    t.Log(modelOutput)
}
