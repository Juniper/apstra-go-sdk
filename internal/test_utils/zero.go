package testutils

func ZeroOf[T any](example T) T {
	var zero T
	return zero
}
