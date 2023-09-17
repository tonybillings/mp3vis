package gfx

// ConcurrentFloater Implementation by external packages should be thread-safe
type ConcurrentFloater interface {
	Float() float32
}
