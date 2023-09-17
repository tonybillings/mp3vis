package gfx

import (
	"context"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"runtime"
	"time"
)

func initWindow(width, height int) (*glfw.Window, error) {
	if err := glfw.Init(); err != nil {
		return nil, err
	}

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "mp3vis", nil, nil)
	if err != nil {
		return nil, err
	}
	window.MakeContextCurrent()

	if err = gl.Init(); err != nil {
		return nil, err
	}

	return window, nil
}

func clearScreen() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func initWindowObjects(window *glfw.Window) {
	for _, o := range windowObjects {
		o.Init(window)
	}
}

func updateWindowObjects(deltaTime int64) {
	for _, o := range windowObjects {
		o.Update(deltaTime)
	}
}

func drawWindowObjects(deltaTime int64) {
	for _, o := range windowObjects {
		o.Draw(deltaTime)
	}
}

func closeWindowObjects() {
	for _, o := range windowObjects {
		o.Close()
	}
}

func Run(ctx context.Context, cancelFunc context.CancelFunc, windowWidth int, windowHeight int) {
	runtime.LockOSThread()

	window, err := initWindow(windowWidth, windowHeight)
	if err != nil {
		fmt.Printf("error initializing OpenGL window: %v\n", err)
		cancelFunc()
	}
	defer glfw.Terminate()

	err = initShaders()
	if err != nil {
		fmt.Printf("error initializing shaders: %v\n", err)
		cancelFunc()
	}

	initWindowObjects(window)

	now := time.Now().UnixMilli()
	lastTick := now
	deltaTime := now

	for !window.ShouldClose() {
		select {
		case <-ctx.Done():
			closeWindowObjects()
			window.SetShouldClose(true)
			return
		default:
		}

		now = time.Now().UnixMilli()
		deltaTime = now - lastTick
		lastTick = now

		clearScreen()

		updateWindowObjects(deltaTime)
		drawWindowObjects(deltaTime)

		window.SwapBuffers()
		glfw.PollEvents()
	}

	cancelFunc()
}
