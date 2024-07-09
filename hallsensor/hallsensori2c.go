package hallsensor

import (
    "math"
    "github.com/d2r2/go-i2c"
    "github.com/d2r2/go-logger"
)


type HallSensorI2C struct {
    i2cChannel *i2c.I2C
    offset int
    prevData int
    resetCount int
}

func NewI2C(address byte, line int) (*HallSensorI2C, error) {
    logger.ChangePackageLogLevel("i2c", logger.FatalLevel)

    i2cChannel, err := i2c.NewI2C(address, line)
    if err != nil {
        return nil, err
    }

    hs := &HallSensorI2C{i2cChannel, 0, 0, 0}

    hs.offset, err = hs.read()
    if err != nil {
        return nil, err
    }
    hs.prevData = hs.offset

    return hs, nil
}

func (hs *HallSensorI2C) read() (int, error) {
    highByte, err := hs.i2cChannel.ReadRegU8(0x00)
    if err != nil {
        return -1, err
    }
    lowByte, err := hs.i2cChannel.ReadRegU8(0x01)
    if err != nil {
        return -1, err
    }

    value := (int(highByte) << 4) | (int(lowByte) >> 4)

    return value, nil
}

func (hs *HallSensorI2C) Read() (float64, error) {
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

    return -float64(output)*STEP_TO_MM, nil
}

func (hs *HallSensorI2C) Close() {
    hs.i2cChannel.Close()
}
