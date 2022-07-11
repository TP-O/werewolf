package itf

type IFactory[K any, V any] interface {
	Create(key K) V
}
