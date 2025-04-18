package utils

func Safe[T any](p *T) interface{} {
	if p == nil {
		return "<nil>"
	}
	return *p
}
