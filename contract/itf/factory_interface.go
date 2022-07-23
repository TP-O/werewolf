package itf

type IFactory[K any, P any, V any] interface {
	Create(key K, payload P) V
}
