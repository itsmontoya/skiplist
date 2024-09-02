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
		switch key.Compare(e.Key) {
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

func (l *level[K]) Insert(seekIndex int, e Entry[K, int]) (index int, err error) {
	if index, err = l.layer.Insert(seekIndex, e); err != nil {
		return
	}

	cur := l.Cursor()
	defer cur.Close()

	iteratingIndex := index + 1
	e, ok := cur.Seek(iteratingIndex)
	for ok {
		e.Value++
		if err = l.Slice.Set(iteratingIndex, e); err != nil {
			return
		}

		e, ok = cur.Next()
		iteratingIndex++
	}

	return
}
