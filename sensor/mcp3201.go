package sensor

import (
	"github.com/roboticslab-uc3m/goraspio/digitalio"
    "github.com/roboticslab-uc3m/goraspio/utils"
)

type Mcp3201 struct {
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

func NewMcp3201(vRef float64, chipSelectPinNo int, voltageInit float64) (*Mcp3201, error) {
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

    return &Mcp3201{
        spi: spi,
        vRef: vRef,
        kf: kf,
        mf: mf,
        voltageInit: voltageInit,
    }, nil
}

func (m *Mcp3201) read() (float64, error) {
    // read bytes
    data, err := m.spi.Read()
    if err != nil {
        return 0.0, err
    }
    
    // convert to voltage
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    voltage := (float64(value) / 4095) * m.vRef

    return voltage, nil
}

func (m *Mcp3201) Read() (VoltageReading, error) {
    voltage, err := m.read()
    if err != nil {
        return VoltageReading{}, err
    }

    // filtering
    voltageMed := m.mf.Compute(voltage)
    voltageFilt := m.kf.Compute(voltageMed)

    return VoltageReading{
        Voltage: voltage,
        VoltageFilt: voltageFilt,
        VoltageRel: (voltage-m.voltageInit) / m.voltageInit * 100,
        VoltageFiltRel: (voltageFilt-m.voltageInit) / m.voltageInit *100,
    }, nil
}

func (m *Mcp3201) Close() error {
    m.spi.Close()

    return nil
}
