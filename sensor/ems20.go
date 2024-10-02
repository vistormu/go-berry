package sensor

import (
	"github.com/roboticslab-uc3m/goraspio/digitalio"
    "github.com/roboticslab-uc3m/goraspio/utils"
)

type Ems20 struct {
    spi digitalio.Spi
    kf *utils.KalmanFilter
    mf *utils.MedianFilter
}

func NewEms20(chipSelectPinNo int) (*Ems20, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo)
    if err != nil {
        return nil, err
    }

    var processVariance float64 = 0.05
    var measurementVariance float64 = 20
    var initialErrorCovariance float64 = 1.0
    kf := utils.NewKalmanFilter(
        processVariance,
        measurementVariance,
        initialErrorCovariance,
    )

    mf := utils.NewMedianFilter(5)
    
    lc :=  &Ems20{
        spi: spi,
        kf: kf,
        mf: mf,
    }
    loadInit, err := lc.read()
    if err != nil {
        return nil, err
    }

    kf.SetInitialEstimate(loadInit)

    return lc, nil
}

func (lc *Ems20) read() (float64, error) {
    // read bytes
    data, err := lc.spi.Read()
    if err != nil {
        return -1, err
    }

    // convert to N
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    load := (float64(value) / 4095) * 50

    return load, nil
}

func (lc *Ems20) Read() (float64, float64, error) {
    load, err := lc.read()
    if err != nil {
        return load, load, nil
    }

    loadMed := lc.mf.Compute(load)
    loadFilt := lc.kf.Compute(loadMed)

    return load, loadFilt, nil
}

func (s *Ems20) Close() error {
    s.spi.Close()

    return nil
}
