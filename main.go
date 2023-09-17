package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"tonysoft.com/mp3vis/audio"
	"tonysoft.com/mp3vis/gfx"
)

func getCliArgs() (mp3Filename, bgFilename string, windowWidth, windowHeight int) {
	if len(os.Args) < 2 {
		panic("must pass in the path to the music file (*.mp3)")
	}

	if len(os.Args) < 3 {
		panic("must pass in the path to the background image (*.png)")
	}

	if len(os.Args) < 4 {
		panic("must pass in the desired window width")
	}

	if len(os.Args) < 5 {
		panic("must pass in the desired window height")
	}

	width, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic("error converting the third CLI argument (window width) to int")
	}

	height, err := strconv.Atoi(os.Args[4])
	if err != nil {
		panic("error converting the fourth CLI argument (window height) to int")
	}

	return os.Args[1], os.Args[2], width, height
}

func waitForInterruptSignal(ctx context.Context, cancelFunc context.CancelFunc) {
	sigIntChan := make(chan os.Signal, 1)
	signal.Notify(sigIntChan, syscall.SIGINT)

	select {
	case <-ctx.Done():
		return
	case <-sigIntChan:
		cancelFunc()
		return
	}
}

func main() {
	mp3Filename, bgFilename, windowWidth, windowHeight := getCliArgs()

	bassPulse := audio.AddBandPeakRegister(25, 80)

	background := gfx.NewPulsingImage(bassPulse, bgFilename)
	circle := gfx.NewPulsingCircle(bassPulse)
	lineMiddle := gfx.NewPulsingLine(bassPulse, 20, 0)
	lineTop := gfx.NewPulsingLine(bassPulse, 20, .5)
	lineBottom := gfx.NewPulsingLine(bassPulse, 20, -.5)

	gfx.AddWindowObject(background)
	gfx.AddWindowObject(circle)
	gfx.AddWindowObject(lineMiddle)
	gfx.AddWindowObject(lineTop)
	gfx.AddWindowObject(lineBottom)

	ctx, cancelFunc := context.WithCancel(context.Background())

	go waitForInterruptSignal(ctx, cancelFunc)
	go audio.Play(ctx, cancelFunc, mp3Filename)
	gfx.Run(ctx, cancelFunc, windowWidth, windowHeight)
}
