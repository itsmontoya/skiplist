package bssl

type Key[K any] interface {
	Compare(K) int
}
