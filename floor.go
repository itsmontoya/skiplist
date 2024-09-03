package skiplist

import (
	"path"
)

func newFloor[K Key, V any](fullPath string) (out *floor[K, V], err error) {
	var f floor[K, V]
	filepath := path.Join(fullPath, "floor.bat")
	if f.layer, err = newLayer[K, V](filepath); err != nil {
		return
	}

	out = &f
	return
}

type floor[K Key, V any] struct {
	*layer[K, V]
}

func (f *floor[K, V]) GetMatch(seekIndex int, key K) (value V, ok bool) {
	cur := f.Cursor()
	defer cur.Close()

	var e Entry[K, V]
	e, ok = cur.Seek(seekIndex)
	for ok {
		switch e.Key.Compare(key) {
		case -1:
			e, ok = cur.Next()
		case 1:
			ok = false
			return
		case 0:
			return e.Value, true
		}
	}

	return
}
