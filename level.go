package bssl

import (
	"fmt"
	"path"

	"github.com/itsmontoya/mappedslice"
)

func newLevel[K Key[any]](fullPath string, n int) (out *level[K], err error) {
	var l level[K]
	filename := fmt.Sprintf("level_%d.bat", n)
	filepath := path.Join(fullPath, filename)
	if l.Slice, err = mappedslice.New[Entry[K, int]](filepath, 32); err != nil {
		return
	}

	l.level = n
	out = &l
	return
}

type level[K Key[any]] struct {
	level int
	*mappedslice.Slice[Entry[K, int]]
}

func (l *level[K]) getIndex(seekIndex int, key K) (index int) {
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
