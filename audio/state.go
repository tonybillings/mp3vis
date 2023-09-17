package audio

import "math"

const (
	epsilon = 0.0001
)

var (
	bpRegisters []*BandPeakRegister
)

func AddBandPeakRegister(freqBegin, freqEnd float32) *BandPeakRegister {
	for _, reg := range bpRegisters {
		if math.Abs(float64(reg.freqBegin-freqBegin)) < epsilon && math.Abs(float64(reg.freqEnd-freqEnd)) < epsilon {
			return reg
		}
	}

	r := newBandPeakRegister(freqBegin, freqEnd)
	bpRegisters = append(bpRegisters, r)
	return r
}
