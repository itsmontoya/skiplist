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
	count int
	*mappedslice.Slice[Entry[K, int]]
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

func (l *level[K]) Insert(key K) (index int, err error) {
	var match bool
	if index, match = l.getIndex(key); match {
		err = fmt.Errorf("key of <%v> is already in level", key)
		return
	}

	var e Entry[K, int]
	e.Key = key
	e.Value = index

	l.Slice.InsertAt(index, e)
	l.count++
	return
}

// TODO: Improve this lookup with binary search
func (l *level[K]) getIndex(key K) (index int, ok bool) {
	cur := l.Cursor()
	defer cur.Close()
	e, ok := cur.Seek(index)
	for ok {
		switch key.Compare(e.Key) {
		case -1:
			index = e.Value
			e, ok = cur.Next()
			index++
		case 1:
			return index, false
		case 0:
			return index, true
		}
	}

	return
}
