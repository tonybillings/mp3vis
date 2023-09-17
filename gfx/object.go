package gfx

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"image"
	"image/draw"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"unsafe"
)

var (
	sizeOfFloat32 = int(unsafe.Sizeof(float32(0)))
)

type PulsingSourceHandle *float32

type WindowObject interface {
	Init(*glfw.Window)
	Update(int64)
	Draw(int64)
	Close()
}

type WindowObjectBase struct {
	window *glfw.Window
	vao    uint32
	vbo    []uint32
	shader uint32
	verts  []float32
}

type PulsingObjectBase struct {
	pulseSource ConcurrentFloater
	pulseScalar float32
}

type PulsingCircle struct {
	WindowObjectBase
	PulsingObjectBase
	scaleUniformPos      int32
	colorScaleUniformPos int32
}

type PulsingLine struct {
	WindowObjectBase
	PulsingObjectBase
	thickness           int32
	thicknessUniformPos int32
	stepUniformPos      int32
	vertexCount         int32
	centerY             float32
}

type PulsingImage struct {
	WindowObjectBase
	PulsingObjectBase
	filename          string
	texture           uint32
	rotAngle          float32
	textureUniformPos int32
	angleUniformPos   int32
}

func (o *WindowObjectBase) Update(_ int64) {
	// Default implementation
}

func (o *WindowObjectBase) Close() {
	gl.BindVertexArray(o.vao)
	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
	gl.DeleteBuffers(int32(len(o.vbo)), &o.vbo[0])
	gl.DeleteVertexArrays(1, &o.vao)
}

func (o *PulsingObjectBase) Update() {
	o.pulseScalar = o.pulseSource.Float()
}

func (c *PulsingCircle) initVertices() {
	c.verts = make([]float32, 2)
	for i := 0; i <= 100; i++ {
		theta := float64(i) * 2.0 * math.Pi / float64(100)
		dx := float32(math.Cos(theta)) * 0.5
		dy := float32(math.Sin(theta)) * 0.5
		c.verts = append(c.verts, dx, dy)
	}
}

func (c *PulsingCircle) initGl() {
	c.shader = shaders[pulsingCircleShader]
	c.scaleUniformPos = gl.GetUniformLocation(c.shader, gl.Str("scale\x00"))
	c.colorScaleUniformPos = gl.GetUniformLocation(c.shader, gl.Str("colorScale\x00"))
	c.vbo = make([]uint32, 1)

	gl.GenVertexArrays(1, &c.vao)
	gl.GenBuffers(1, &c.vbo[0])

	gl.BindVertexArray(c.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(c.verts)*sizeOfFloat32, gl.Ptr(c.verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (l *PulsingLine) initVertices() {
	l.vertexCount = 200
	l.verts = make([]float32, l.vertexCount*2)
	step := 2.0 / float32(l.vertexCount)

	for i := int32(0); i < l.vertexCount; i++ {
		l.verts[i*2] = -1 + float32(i)*step
		l.verts[i*2+1] = l.centerY
	}
}

func (l *PulsingLine) initGl() {
	l.shader = shaders[pulsingLineShader]
	l.thicknessUniformPos = gl.GetUniformLocation(l.shader, gl.Str("thickness\x00"))
	l.stepUniformPos = gl.GetUniformLocation(l.shader, gl.Str("step\x00"))

	_, height := l.window.GetSize()
	pixelHeight := 2.0 / float32(height)

	gl.UseProgram(l.shader)
	gl.Uniform1f(l.stepUniformPos, pixelHeight)
	gl.UseProgram(0)

	l.vbo = make([]uint32, 1)

	gl.GenVertexArrays(1, &l.vao)
	gl.GenBuffers(1, &l.vbo[0])

	gl.BindVertexArray(l.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo[0])

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, 0, nil)

	gl.BufferData(gl.ARRAY_BUFFER, len(l.verts)*sizeOfFloat32, gl.Ptr(l.verts), gl.DYNAMIC_DRAW)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (i *PulsingImage) initVertices() {
	i.verts = []float32{
		-1.0, 1.0, 0.0, 0.0,
		1.0, 1.0, 1.0, 0.0,
		1.0, -1.0, 1.0, 1.0,
		-1.0, -1.0, 0.0, 1.0,
	}
}

func (i *PulsingImage) initTexture() {
	imgFile, err := os.Open(i.filename)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		e := imgFile.Close()
		if e != nil {
			panic(e)
		}
	}(imgFile)

	img, _, err := image.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	gl.GenTextures(1, &i.texture)
	gl.BindTexture(gl.TEXTURE_2D, i.texture)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix),
	)

	gl.GenerateMipmap(gl.TEXTURE_2D)
}

func (i *PulsingImage) initGl() {
	i.shader = shaders[pulsingImageShader]
	i.angleUniformPos = gl.GetUniformLocation(i.shader, gl.Str("angle\x00"))
	i.textureUniformPos = gl.GetUniformLocation(i.shader, gl.Str("backgroundTexture\x00"))
	i.vbo = make([]uint32, 1)

	gl.GenVertexArrays(1, &i.vao)
	gl.GenBuffers(1, &i.vbo[0])

	gl.BindVertexArray(i.vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, i.vbo[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(i.verts)*sizeOfFloat32, gl.Ptr(i.verts), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointerWithOffset(0, 2, gl.FLOAT, false, 4*4, 0)
	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointerWithOffset(1, 2, gl.FLOAT, false, 4*4, 2*4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (c *PulsingCircle) Init(window *glfw.Window) {
	c.window = window
	c.initVertices()
	c.initGl()
}

func (l *PulsingLine) Init(window *glfw.Window) {
	l.window = window
	l.initVertices()
	l.initGl()
}

func (i *PulsingImage) Init(window *glfw.Window) {
	i.window = window
	i.initVertices()
	i.initTexture()
	i.initGl()
}

func (c *PulsingCircle) Update(deltaTime int64) {
	c.PulsingObjectBase.Update()
}

func (l *PulsingLine) Update(deltaTime int64) {
	l.PulsingObjectBase.Update()

	maxOffset := float32(0.25)
	noise := make([]float32, l.vertexCount)
	for i := int32(0); i < l.vertexCount; i++ {
		noise[i] = float32((rand.Float64() * 2.0) - 1.0)
		l.verts[i*2+1] = l.centerY + (noise[i] * l.pulseScalar * maxOffset)
	}
}

func (i *PulsingImage) Update(deltaTime int64) {
	i.PulsingObjectBase.Update()
}

func (c *PulsingCircle) Draw(deltaTime int64) {
	gl.UseProgram(c.shader)

	gl.BindVertexArray(c.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, c.vbo[0])

	gl.Uniform1f(c.scaleUniformPos, c.pulseScalar)
	gl.Uniform1f(c.colorScaleUniformPos, c.pulseScalar)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, int32(len(c.verts)/2))

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func (l *PulsingLine) Draw(deltaTime int64) {
	gl.UseProgram(l.shader)
	gl.Uniform1f(l.thicknessUniformPos, float32(l.thickness))

	gl.BindVertexArray(l.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, l.vbo[0])

	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(l.verts)*sizeOfFloat32, gl.Ptr(l.verts))

	gl.DrawArrays(gl.LINE_STRIP, 0, l.vertexCount)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func (i *PulsingImage) Draw(deltaTime int64) {
	gl.UseProgram(i.shader)

	gl.BindVertexArray(i.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, i.vbo[0])

	gl.Uniform1f(i.angleUniformPos, 0.5*i.pulseScalar)
	gl.Uniform1i(i.textureUniformPos, 0)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, i.texture)

	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)

	gl.BindVertexArray(0)
	gl.UseProgram(0)
}

func (i *PulsingImage) Close() {
	i.WindowObjectBase.Close()
	gl.DeleteTextures(1, &i.texture)
}

func NewPulsingCircle(pulseSource ConcurrentFloater) *PulsingCircle {
	return &PulsingCircle{
		PulsingObjectBase: PulsingObjectBase{
			pulseSource: pulseSource,
		},
	}
}

func NewPulsingLine(pulseSource ConcurrentFloater, thickness int32, verticalPosition float32) *PulsingLine {
	return &PulsingLine{
		PulsingObjectBase: PulsingObjectBase{
			pulseSource: pulseSource,
		},
		thickness: thickness,
		centerY:   verticalPosition,
	}
}

func NewPulsingImage(pulseSource ConcurrentFloater, filename string) *PulsingImage {
	return &PulsingImage{
		PulsingObjectBase: PulsingObjectBase{
			pulseSource: pulseSource,
		},
		filename: filename,
	}
}
