package bssl

type Entry[K Key[any], V any] struct {
	Key   K
	Value V
}
