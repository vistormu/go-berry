package loadcell

import (
    "goraspio/digitalio"
)

type LoadCell struct {
    spi digitalio.Spi
}

func New(chipSelectPinNo int) (LoadCell, error) {
    spi, err := digitalio.NewSpi(chipSelectPinNo)
    if err != nil {
        return LoadCell{}, err
    }

    return LoadCell{spi}, nil
}

func (lc LoadCell) Read() (float32, error) {
    // read bytes
    data, err := lc.spi.Read()
    if err != nil {
        return -1, err
    }

    // convert to kg
    value := ((int(data[0]) & 0x1F) << 7) | (int(data[1]) >> 1)
    load := (float32(value) / 4095) * 5

    return load, nil
}

func (ld LoadCell) Close() {
    ld.spi.Close()
}
