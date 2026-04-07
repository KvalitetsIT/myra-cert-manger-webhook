package client

// Implement this interface in order to provide the solver another DNS client
type Client[T any] interface {
	OnDelete(record T) (T, error)
	OnAdd(record T) (T, error)
}
