package model

type Requestable[T, S any] interface {
	Make(S) (*T, error)
}

func MakeRequest[T Requestable[T, S], S any](request S) (*T, error) {
	req := new(T)

	result, err := (*req).Make(request)

	return result, err
}
