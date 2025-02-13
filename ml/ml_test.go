package ml

import (
    "fmt"
    "testing"
    "time"
)

func TestMain(t *testing.T) {
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

func TestServer(t *testing.T) {
    model, _ := NewRemoteModel(
        "145.94.127.212",
        8080,
        "predict",
        "output",
    )

    start := time.Now()
    modelOutput, err := model.Compute(struct {
        S0 float64 `json:"s0"`
        S1 float64 `json:"s1"`
        S2 float64 `json:"s2"`
        S3 float64 `json:"s3"`
    }{
        0.0, 0.0, 0.0, 0.0,
    })
    if err != nil {
        t.Fatal(err)
    }
    
    t.Log(time.Since(start))
    t.Log(modelOutput)
}
