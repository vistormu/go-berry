package hallsensor

import (
    "math"
    "goraspio/digitalio"
)

const (
    STEP_TO_MM = 0.000488
)

type HallSensorSpi struct {
    spi digitalio.Spi 
    offset int
    prevData int
    resetCount int
}

func NewSpi(chipSelectNo int) (*HallSensorSpi, error) {
    spi, err := digitalio.NewSpi(chipSelectNo)
    if err != nil {
        return nil, err
    }

    hs := &HallSensorSpi{spi, 0, 0, 0}

    hs.offset, err = hs.read()
    if err != nil {
        return nil, err
    }
    hs.prevData = hs.offset

    return hs, nil
}

func (hs *HallSensorSpi) read() (int, error) {
    data, err := hs.spi.Read()
    if err != nil {
        return -1, err
    }

    value := (int(data[0]) << 4) | (int(data[1]) >> 4)

    return value, nil
}

func (hs *HallSensorSpi) Read() (float64, error) {
    data, err := hs.read()
    if err != nil {
        return -1, err
    }

    diff := float64(data - hs.prevData)
    change := float64(MAX_VALUE)*(1-RESET_THRESH)
    if diff < 0 && math.Abs(diff) > change {
        hs.resetCount++
    }
    if diff > 0 && math.Abs(diff) > change {
        hs.resetCount--
    }

    hs.prevData = data
    output := data - hs.offset + hs.resetCount * (MAX_VALUE + 1)

    return float64(output)*STEP_TO_MM, nil
}

func (hs *HallSensorSpi) Close() {
    hs.spi.Close()
}
