package gfx

// ConcurrentFloater Implementation by external packages should be thread-safe
type ConcurrentFloater interface {
	Float() float32
}

type DefaultFloater struct{}

// Float This implementation does not need to be thread-safe
func (f *DefaultFloater) Float() float32 {
	return 0
}
