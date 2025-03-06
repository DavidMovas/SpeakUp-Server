package requests

type requestable[T, S any] interface {
	make(S) (*T, error)
}

func MakeRequest[T requestable[T, S], S any](request S) (*T, error) {
	req := new(T)

	result, err := (*req).make(request)

	return result, err
}
