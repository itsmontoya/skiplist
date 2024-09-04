package skiplist

func makeEntry[K Key, V any](key *K, val *V) (e Entry[K, V]) {
	e.Key = *key
	e.Value = *val
	return
}

type Entry[K Key, V any] struct {
	Key   K
	Value V
}
