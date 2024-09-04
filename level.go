package skiplist

import (
	"fmt"
	"path"
)

func newLevel[K Key](fullPath string, n int) (out *level[K], err error) {
	var l level[K]
	filename := fmt.Sprintf("level_%d.bat", n)
	filepath := path.Join(fullPath, filename)
	if l.layer, err = newLayer[K, int](filepath, 32); err != nil {
		return
	}

	l.level = n
	out = &l
	return
}

type level[K Key] struct {
	level int
	*layer[K, int]
}

func (l *level[K]) GetSeekIndex(seekIndex int, key *K) (index int) {
	for e, ok := l.Slice.Get(seekIndex); ok; e, ok = l.Slice.Get(seekIndex) {
		switch e.Key.Compare(key) {
		case -1:
			index = e.Value
			seekIndex++
		case 1:
			return index
		case 0:
			return e.Value
		}
	}

	return
}

func (l *level[K]) IterateAfter(index int, key *K, fn func(index int, e Entry[K, int])) (seekIndex int) {
	for e, ok := l.Get(index); ok; e, ok = l.Get(index) {
		switch {
		case e.Key.Compare(key) < 1:
			seekIndex = index
		default:
			fn(index, e)
		}

		index++
	}

	return
}
