package model

import (
    "fmt"
    "errors"
    ort "github.com/yalue/onnxruntime_go"
)

type ModelType int
const (
    Transformer ModelType = iota
)

var modelTypeToPath = map[ModelType]string{
    Transformer: "model/models/TransformerRegressor.onnx",
}


type Model struct {
    contextLength int
    session *ort.AdvancedSession
    inputTensor *ort.Tensor[float64]
    outputTensor *ort.Tensor[float64]
}

func New(modelType ModelType, contextLength int) (*Model, error) {
    path, ok := modelTypeToPath[modelType]
    if !ok {
        return nil, errors.New("Unknown model type")
    }

    // environment
    ort.SetSharedLibraryPath("model/onnxruntime/lib/libonnxruntime.so")

    err := ort.InitializeEnvironment()
    if err != nil {
        return nil, err
    }

    // tensors
    input := make([]float64, contextLength)
    inputTensor, err := ort.NewTensor(ort.NewShape(1, int64(contextLength)), input)
    if err != nil {
        return nil, err
    }

    outputTensor, err := ort.NewEmptyTensor[float64](ort.NewShape(1, 1))
    if err != nil {
        return nil, err
    }

    // session
    session, err := ort.NewAdvancedSession(
        path,
        []string{"input"},
        []string{"output"},
        []ort.ArbitraryTensor{inputTensor},
		[]ort.ArbitraryTensor{outputTensor},
        nil,
    )
    if err != nil {
        return nil, err
    }

    return &Model{
        contextLength,
        session,
        inputTensor,
        outputTensor,
    }, nil
}

func (m *Model) Compute(input []float64) (float64, error) {
    if len(input) != m.contextLength {
        return 0.0, fmt.Errorf("Input dimension should match given context length. Got %d and want %d", len(input), m.contextLength)
    }

    copy(m.inputTensor.GetData(), input)

    err := m.session.Run()
    if err != nil {
        return 0.0, err
    }

    output := m.outputTensor.GetData()[0]

    return output, nil
}


func (m *Model) Close() {
    m.inputTensor.Destroy()
    m.outputTensor.Destroy()
    m.session.Destroy()
    ort.DestroyEnvironment()
} 

