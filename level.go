package skiplist

import (
	"fmt"
	"path"
)

func newLevel[K Key](fullPath string, n int) (out *level[K], err error) {
	var l level[K]
	filename := fmt.Sprintf("level_%d.bat", n)
	filepath := path.Join(fullPath, filename)
	if l.layer, err = newLayer[K, int](filepath); err != nil {
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

func (l *level[K]) GetSeekIndex(seekIndex int, key K) (index int) {
	cur := l.Cursor()
	defer cur.Close()

	e, ok := cur.Seek(seekIndex)
	for ok {
		switch e.Key.Compare(key) {
		case -1:
			index = e.Value
			e, ok = cur.Next()
		case 1:
			return index
		case 0:
			return e.Value
		}
	}

	return
}

func (l *level[K]) IterateAfter(index int, key K, fn func(index int, e Entry[K, int])) (seekIndex int) {
	cur := l.Cursor()
	defer cur.Close()

	e, ok := cur.Seek(index)
	for ok {
		if e.Key.Compare(key) < 1 {
			seekIndex = e.Value
		} else {
			fn(index, e)
		}

		e, ok = cur.Next()
		index++
	}

	return
}
