package contract

type Factory[K any, T any] interface {
	Create(key K) (T, error)
}
