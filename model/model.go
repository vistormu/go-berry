package model

import (
    "os"
    "fmt"
    ort "github.com/yalue/onnxruntime_go"
)

type Model struct {
    contextLength int
    session *ort.AdvancedSession
    inputTensor *ort.Tensor[float32]
    outputTensor *ort.Tensor[float32]
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil // Path exists
    }
    if os.IsNotExist(err) {
        return false, nil // Path does not exist
    }
    return false, err // Some other error
}

func New(modelName string, contextLength int) (*Model, error) {
    path := "model/models/" + modelName + ".onnx"
    // path := "models/" + modelName + ".onnx"
    _, err := os.Stat(path)
    if os.IsNotExist(err) {
        return nil, err
    }

    // environment
    ort.SetSharedLibraryPath("model/onnxruntime/lib/libonnxruntime.so")
    // ort.SetSharedLibraryPath("onnxruntime/lib/libonnxruntime.so")

    err = ort.InitializeEnvironment()
    if err != nil {
        return nil, err
    }

    // tensors
    input := make([]float32, contextLength)
    inputTensor, err := ort.NewTensor(ort.NewShape(1, int64(contextLength)), input)
    if err != nil {
        return nil, err
    }

    outputTensor, err := ort.NewEmptyTensor[float32](ort.NewShape(1))
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

    modelInput := make([]float32, len(input))
    for i, v := range input {
        modelInput[i] = float32(v)
    }

    copy(m.inputTensor.GetData(), modelInput)

    err := m.session.Run()
    if err != nil {
        return 0.0, err
    }

    output := m.outputTensor.GetData()[0]

    return float64(output), nil
}


func (m *Model) Close() {
    m.inputTensor.Destroy()
    m.outputTensor.Destroy()
    m.session.Destroy()
    ort.DestroyEnvironment()
} 

