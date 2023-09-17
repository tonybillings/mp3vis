package gfx

import (
	_ "embed"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

const (
	pulsingCircleShader = "pulsing_circle"
	pulsingLineShader   = "pulsing_line"
	pulsingImageShader  = "pulsing_image"
)

var (
	//go:embed shaders/pulsing_circle.vs
	pulsingCircleShaderVs string

	//go:embed shaders/pulsing_circle.fs
	pulsingCircleShaderFs string

	//go:embed shaders/pulsing_line.vs
	pulsingLineShaderVs string

	//go:embed shaders/pulsing_line.fs
	pulsingLineShaderFs string

	//go:embed shaders/pulsing_image.vs
	pulsingImageShaderVs string

	//go:embed shaders/pulsing_image.fs
	pulsingImageShaderFs string
)

var (
	shaders map[string]uint32
)

func initShaders() error {
	circleProg, err := createShaderProgram(pulsingCircleShaderVs, pulsingCircleShaderFs)
	if err != nil {
		return err
	}

	lineProg, err := createShaderProgram(pulsingLineShaderVs, pulsingLineShaderFs)
	if err != nil {
		return err
	}

	imageProg, err := createShaderProgram(pulsingImageShaderVs, pulsingImageShaderFs)
	if err != nil {
		return err
	}

	shaders = make(map[string]uint32)
	shaders[pulsingCircleShader] = circleProg
	shaders[pulsingLineShader] = lineProg
	shaders[pulsingImageShader] = imageProg

	return nil
}

func checkShaderError(shader uint32) error {
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to compile shader: %s", log)
	}
	return nil
}

func checkProgramError(program uint32) error {
	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return fmt.Errorf("failed to link program: %s", log)
	}
	return nil
}

func createShaderProgram(vertexSource, fragmentSource string) (uint32, error) {
	vertexShader := gl.CreateShader(gl.VERTEX_SHADER)
	cstr, free := gl.Strs(vertexSource + "\x00")
	gl.ShaderSource(vertexShader, 1, cstr, nil)
	free()
	gl.CompileShader(vertexShader)
	if err := checkShaderError(vertexShader); err != nil {
		return 0, err
	}

	fragmentShader := gl.CreateShader(gl.FRAGMENT_SHADER)
	cstr, free = gl.Strs(fragmentSource + "\x00")
	gl.ShaderSource(fragmentShader, 1, cstr, nil)
	free()
	gl.CompileShader(fragmentShader)
	if err := checkShaderError(fragmentShader); err != nil {
		return 0, err
	}

	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertexShader)
	gl.AttachShader(shaderProgram, fragmentShader)
	gl.LinkProgram(shaderProgram)
	if err := checkProgramError(shaderProgram); err != nil {
		return 0, err
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return shaderProgram, nil
}
