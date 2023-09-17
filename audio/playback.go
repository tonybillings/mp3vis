package audio

import (
	"context"
	"fmt"
	"github.com/hajimehoshi/go-mp3"
	"github.com/mesilliac/pulse-simple"
	"gonum.org/v1/gonum/dsp/fourier"
	"math/cmplx"
	"os"
)

const (
	fftSize = 4096
)

func updateBandPeakRegisters(freq float32, coef complex128) {
	for _, r := range bpRegisters {
		if freq > r.freqBegin && freq < r.freqEnd {
			amplitude := float32(cmplx.Abs(coef)) / float32(fftSize)
			if amplitude > r.ampMax {
				r.ampMax = amplitude
			}
			if amplitude > r.curAmpMax {
				r.curAmpMax = amplitude
			}
		}
	}

}

func setBandPeakRegisters() {
	for _, r := range bpRegisters {
		r.set(r.curAmpMax / r.ampMax)
		r.curAmpMax = 0
	}
}

func Play(ctx context.Context, cancelFunc context.CancelFunc, mp3Filename string) {
	file, err := os.Open(mp3Filename)
	if err != nil {
		fmt.Printf("error opening music file: %v\n", err)
		cancelFunc()
	}
	defer file.Close()

	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		fmt.Printf("error creating new mp3 decoder: %v\n", err)
		cancelFunc()
	}

	ss := pulse.SampleSpec{Format: pulse.SAMPLE_S16LE, Rate: uint32(decoder.SampleRate()), Channels: 2}
	playback, err := pulse.Playback("mp3vis", "mp3vis", &ss)
	if err != nil {
		fmt.Printf("error creating pulse playback stream: %v\n", err)
		cancelFunc()
	}

	defer playback.Free()
	defer playback.Drain()

	fft := fourier.NewFFT(fftSize)
	freqDomain := make([]complex128, fftSize/2+1)
	sampleBuffer := make([]float64, fftSize)

	freqResolution := float32(decoder.SampleRate()) / float32(fftSize)

	buf := make([]byte, 2*fftSize) // assuming 16-bit samples
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		n, readErr := decoder.Read(buf)
		if err != nil {
			fmt.Printf("error reading buffer: %v\n", readErr)
			cancelFunc()
		}
		if n == 0 {
			break
		}

		for i := 0; i < n; i += 4 {
			val := int16(buf[i]) | (int16(buf[i+1]) << 8) // Use only one channel
			sampleBuffer[i/4] = float64(val)
		}

		_, writeErr := playback.Write(buf[:n])
		if writeErr != nil {
			fmt.Printf("error writing buffer to PulseAudio stream: %v\n", writeErr)
			cancelFunc()
		}

		fft.Coefficients(freqDomain, sampleBuffer)

		for i, coef := range freqDomain {
			freq := float32(i) * freqResolution
			updateBandPeakRegisters(freq, coef)
		}

		setBandPeakRegisters()
	}

	cancelFunc()
}
