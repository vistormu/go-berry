package ml

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/vistormu/go-berry/errors"
)

type RemoteModel struct {
    url         string
    responseKey string
    client      *http.Client
    inputNames  []string
}

func NewRemoteModel(ip string, port int, endpoint string, responseKey string, inputNames []string) (*RemoteModel, error) {
    client := &http.Client{
        Timeout: 10 * time.Second, // adjust based on your expected response times
        Transport: &http.Transport{
            MaxIdleConns:        100,
            IdleConnTimeout:     90 * time.Second,
            DisableCompression:  false,
        },
    }
    
    return &RemoteModel{
        url:         fmt.Sprintf("http://%s:%d/%s", ip, port, endpoint),
        responseKey: responseKey,
        client:      client,
        inputNames:  inputNames,
    }, nil
}

func (m *RemoteModel) Compute(input []float64) ([]float32, error) {
    data := make(map[string]float64)
    for i, name := range m.inputNames {
        data[name] = input[i]
    }

    jsonData, err := json.Marshal(data)
    if err != nil {
        return nil, errors.New(errors.JSON_ENCODE, err)
    }
    
    resp, err := m.client.Post(m.url, "application/json", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, errors.New(errors.MODEL_REQUEST, err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New(errors.STATUS_CODE, resp.StatusCode, resp.Status)
    }

    var result map[string][]float32
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, errors.New(errors.JSON_DECODE, err)
    }

    return result[m.responseKey], nil
}

func (m *RemoteModel) Close() {
    m.client.CloseIdleConnections()
}
