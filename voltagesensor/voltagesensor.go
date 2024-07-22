package voltagesensor

import (
    "goraspio/digitalio"
    "goraspio/utils"
)

type VoltageSensor struct {
    spi digitalio.Spi
    vRef float64
    kf *utils.KalmanFilter
    kfInitialized bool
    prevValue float64
}

func New(vRef float64, chipSelectPinNo int) (*VoltageSensor, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo) 
    if err != nil {
        return nil, err
    }
    
    processVariance := 0.05
    measurementVariance := 30.0
    initialErrorCovariance := 1.0
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
        prevValue: 0.0,
    }, nil
}

func (vs *VoltageSensor) Read() (float64, float64, error) {
    // read bytes
    data, err := vs.spi.Read()
    if err != nil {
        return -1, -1, err
    }
    
    // convert to voltage
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float64(value) / 4095) * vs.vRef

    // filtering
    if !vs.kfInitialized {
        vs.kf.SetInitialEstimate(voltage)
        vs.kfInitialized = true
        vs.prevValue = voltage

        return voltage, voltage, nil
    }
    if voltage > 1.25*vs.prevValue {
        voltage = vs.prevValue
    }

    filteredVoltage := vs.kf.Compute(voltage)

    vs.prevValue = voltage

    return voltage, filteredVoltage, nil
}

func (vs *VoltageSensor) Close() {
    vs.spi.Close()
}
