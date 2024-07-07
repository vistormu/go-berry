package voltagesensor

import (
    "goraspio/digitalio"
    "goraspio/utils"
)

type VoltageSensor struct {
    spi digitalio.Spi
    vRef float32
    kf *utils.KalmanFilter
    kfInitialized bool
}

func New(vRef float32, chipSelectPinNo int) (*VoltageSensor, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo) 
    if err != nil {
        return nil, err
    }
    
    var processVariance float32 = 0.01
    var measurementVariance float32 = 0.1
    var initialErrorCovariance float32 = 1.0
    kf := utils.NewKalmanFilter(
        processVariance,
        measurementVariance,
        initialErrorCovariance,
    )

    return &VoltageSensor{
        spi: spi,
        vRef: vRef,
        kf: kf,
        kfInitialized: false,
    }, nil
}

func (vs *VoltageSensor) Read() (float32, float32, error) {
    // read bytes
    data, err := vs.spi.Read()
    if err != nil {
        return -1, -1, err
    }
    
    // convert to voltage
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float32(value) / 4095) * vs.vRef

    // filtering
    if !vs.kfInitialized {
        vs.kf.SetInitialEstimate(voltage)
        vs.kfInitialized = true

        return voltage, voltage, nil
    }
    filteredVoltage := vs.kf.Compute(voltage)

    return voltage, filteredVoltage, nil
}

func (vs *VoltageSensor) Close() {
    vs.spi.Close()
}
