package skiplist

import (
	"path"
)

func newFloor[K Key, V any](fullPath string, size int64) (out *floor[K, V], err error) {
	var f floor[K, V]
	filepath := path.Join(fullPath, "floor.bat")
	if f.layer, err = newLayer[K, V](filepath, size); err != nil {
		return
	}

	out = &f
	return
}

type floor[K Key, V any] struct {
	*layer[K, V]
}

func (f *floor[K, V]) GetMatch(seekIndex int, key *K) (value V, match bool) {
	for e, ok := f.Get(seekIndex); ok; e, ok = f.Get(seekIndex) {
		switch e.Key.Compare(key) {
		case -1:
			seekIndex++
		case 1:
			match = false
			return
		case 0:
			return e.Value, true
		}
	}

	return
}
