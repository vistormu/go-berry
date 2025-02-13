package ml

import (
    "os"
    "path/filepath"
    _ "embed"
    ort "github.com/yalue/onnxruntime_go"

    "github.com/vistormu/go-berry/errors"
)

//go:embed lib/libonnxruntime.so
var onnxRuntimeData []byte

type LocalModel struct {
    nInputs int
    nOutputs int
    session *ort.AdvancedSession
    inputTensor *ort.Tensor[float32]
    outputTensor *ort.Tensor[float32]
}

func NewLocalModel(modelPath string, nInputs, nOutputs int) (*LocalModel, error) {
    // model path
    _, err := os.Stat(modelPath)
    if os.IsNotExist(err) {
        return nil, errors.New(errors.MODEL_PATH, modelPath, err)
    }
    
    // lib path
    libDir := filepath.Join(os.TempDir(), "lib")
    os.MkdirAll(libDir, 0755)
    libPath := filepath.Join(libDir, "onnxruntime.so")

    err = os.WriteFile(libPath, onnxRuntimeData, 0755)
    if err != nil {
        return nil, err
    }

    ort.SetSharedLibraryPath(libPath)

    err = ort.InitializeEnvironment()
    if err != nil {
        return nil, err
    }

    // tensors
    inputData := make([]float32, nInputs)
    inputShape := ort.NewShape(1, int64(nInputs))
    inputTensor, err := ort.NewTensor(inputShape, inputData)
    if err != nil {
        return nil, err
    }
    
    outputData := make([]float32, nOutputs)
    outputShape := ort.NewShape(int64(nOutputs))
    outputTensor, err := ort.NewTensor(outputShape, outputData)
    if err != nil {
        return nil, err
    }

    // session
    session, err := ort.NewAdvancedSession(
        modelPath,
        []string{"input"},
        []string{"output"},
        []ort.Value{inputTensor},
		[]ort.Value{outputTensor},
        nil,
    )
    if err != nil {
        return nil, err
    }

    return &LocalModel{
        nInputs,
        nOutputs,
        session,
        inputTensor,
        outputTensor,
    }, nil
}

func (m *LocalModel) Compute(input []float64) ([]float32, error) {
    if len(input) != m.nInputs {
        return nil, errors.New(errors.SHAPE_MISMATCH, m.nInputs, len(input))
    }

    modelInput := make([]float32, len(input))
    for i, v := range input {
        modelInput[i] = float32(v)
    }

    copy(m.inputTensor.GetData(), modelInput)

    err := m.session.Run()
    if err != nil {
        return nil, err
    }

    output := m.outputTensor.GetData()

    return output, nil
}

func (m *LocalModel) Close() {
    m.inputTensor.Destroy()
    m.outputTensor.Destroy()
    m.session.Destroy()
    ort.DestroyEnvironment()
} 
