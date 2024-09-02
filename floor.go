package bssl

import (
	"path"

	"github.com/itsmontoya/mappedslice"
)

func newFloor[K Key[any], V any](fullPath string) (out *floor[K, V], err error) {
	var f floor[K, V]
	filepath := path.Join(fullPath, "floor.bat")
	if f.Slice, err = mappedslice.New[Entry[K, V]](filepath, 32); err != nil {
		return
	}

	out = &f
	return
}

type floor[K Key[any], V any] struct {
	*mappedslice.Slice[Entry[K, V]]
}

func (f *floor[K, V]) getMatch(seekIndex int, key K) (value V, ok bool) {
	cur := f.Cursor()
	defer cur.Close()

	var e Entry[K, V]
	e, ok = cur.Seek(seekIndex)
	for ok {
		switch key.Compare(e.Key) {
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
