package voltagesensor

import (
    "goraspio/digitalio"
    "goraspio/utils"
)

type VoltageSensor struct {
    spi digitalio.Spi
    vRef float64
    kf *utils.KalmanFilter
    mf *utils.MedianFilter
    voltageInit float64
}

type VoltageReading struct {
    Voltage float64
    VoltageRel float64
    VoltageFilt float64
    VoltageFiltRel float64
}

func New(vRef float64, chipSelectPinNo int, voltageInit float64) (*VoltageSensor, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo) 
    if err != nil {
        return nil, err
    }
    
    processVariance := 0.05
    measurementVariance := 20.0
    initialErrorCovariance := 1.0
    kf := utils.NewKalmanFilter(
        processVariance,
        measurementVariance,
        initialErrorCovariance,
    )
    kf.SetInitialEstimate(voltageInit)

    mf := utils.NewMedianFilter(5)

    return &VoltageSensor{
        spi: spi,
        vRef: vRef,
        kf: kf,
        mf: mf,
        voltageInit: voltageInit,
    }, nil
}

func (vs *VoltageSensor) read() (float64, error) {
    // read bytes
    data, err := vs.spi.Read()
    if err != nil {
        return 0.0, err
    }
    
    // convert to voltage
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float64(value) / 4095) * vs.vRef

    return voltage, nil
}

func (vs *VoltageSensor) Read() (VoltageReading, error) {
    voltage, err := vs.read()
    if err != nil {
        return VoltageReading{}, err
    }

    // filtering
    voltageMed := vs.mf.Compute(voltage)
    voltageFilt := vs.kf.Compute(voltageMed)

    return VoltageReading{
        Voltage: voltage,
        VoltageFilt: voltageFilt,
        VoltageRel: (voltage-vs.voltageInit) / vs.voltageInit * 100,
        VoltageFiltRel: (voltageFilt-vs.voltageInit) / vs.voltageInit *100,
    }, nil
}

func (vs *VoltageSensor) Close() {
    vs.spi.Close()
}
