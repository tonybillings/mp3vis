package gfx

var (
	windowObjects []WindowObject
)

func AddWindowObject(object WindowObject) {
	windowObjects = append(windowObjects, object)
}
