package loadcell

import (
    "goraspio/digitalio"
    "goraspio/utils"
)

type LoadCell struct {
    spi digitalio.Spi
    kf *utils.KalmanFilter
    kfInitialized bool
}

func New(chipSelectPinNo int) (*LoadCell, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo)
    if err != nil {
        return nil, err
    }

    var processVariance float64 = 0.1
    var measurementVariance float64 = 20
    var initialErrorCovariance float64 = 1.0
    kf := utils.NewKalmanFilter(
        processVariance,
        measurementVariance,
        initialErrorCovariance,
    )

    return &LoadCell{spi, kf, false}, nil
}

func (lc *LoadCell) Read() (float64, float64, error) {
    // read bytes
    data, err := lc.spi.Read()
    if err != nil {
        return -1, -1, err
    }

    // convert to N
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    load := (float64(value) / 4095) * 50

    // filtering
    if !lc.kfInitialized {
        lc.kf.SetInitialEstimate(load)
        lc.kfInitialized = true

        return load, load, nil
    }
    filteredLoad := lc.kf.Compute(load)

    return load, filteredLoad, nil
}

func (ld *LoadCell) Close() {
    ld.spi.Close()
}
