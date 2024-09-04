package skiplist

import (
	"fmt"

	"github.com/itsmontoya/mappedslice"
)

func newLayer[K Key, V any](filepath string, size int64) (out *layer[K, V], err error) {
	var l layer[K, V]
	if l.Slice, err = mappedslice.New[Entry[K, V]](filepath, size); err != nil {
		return
	}

	out = &l
	return
}

type layer[K Key, V any] struct {
	count int
	*mappedslice.Slice[Entry[K, V]]
}

func (l *layer[K, V]) Insert(seekIndex int, e Entry[K, V]) (index int, err error) {
	var match bool
	if index, match = l.getIndex(seekIndex, &e.Key); match {
		err = fmt.Errorf("key of <%v> is already in layer", e.Key)
		return
	}

	err = l.InsertAt(index, e)
	return
}

func (l *layer[K, V]) InsertAt(index int, e Entry[K, V]) (err error) {
	if err = l.insertAt(index, e); err != nil {
		return
	}

	l.count++
	return
}

func (l *layer[K, V]) insertAt(index int, e Entry[K, V]) (err error) {
	if index >= l.Slice.Len() {
		return l.Slice.Append(e)
	}

	return l.Slice.InsertAt(index, e)
}

// TODO: Improve this lookup with binary search
func (l *layer[K, V]) getIndex(seekIndex int, key *K) (index int, ok bool) {
	index = seekIndex
	for e, ok := l.Get(index); ok; e, ok = l.Get(index) {
		switch e.Key.Compare(key) {
		case -1:
			index++
		case 1:
			return index, false
		case 0:
			return index, true
		}
	}

	return
}

func (l *layer[K, V]) printLayer() {
	var i int
	fmt.Print("[")
	l.Slice.ForEach(func(e Entry[K, V]) (end bool) {
		if i > 0 {
			fmt.Print(", ")
		}

		fmt.Printf("%v:%v", e.Key, e.Value)
		i++
		return false
	})
	fmt.Print("]\n")
}
