package audio

import "sync"

type RegisterType interface {
	float32
}

type BaseRegister[T RegisterType] struct {
	value      T
	valueMutex sync.Mutex

	freqBegin float32
	freqEnd   float32
	ampMax    float32
	curAmpMax float32
}

func (r *BaseRegister[T]) Get() T {
	r.valueMutex.Lock()
	val := r.value
	r.valueMutex.Unlock()
	return val
}

func (r *BaseRegister[T]) set(val T) {
	r.valueMutex.Lock()
	r.value = val
	r.valueMutex.Unlock()
}

type BandPeakRegister struct {
	BaseRegister[float32]
}

func (r *BandPeakRegister) Float() float32 {
	return r.Get()
}

func newBandPeakRegister(freqBegin, freqEnd float32) *BandPeakRegister {
	r := &BandPeakRegister{}
	r.freqBegin = freqBegin
	r.freqEnd = freqEnd
	return r
}
