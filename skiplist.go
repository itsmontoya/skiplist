package skiplist

import (
	"errors"
	"fmt"
	"os"
	"path"
)

var (
	ErrNotFound = errors.New("entry not found")
)

func New[K Key, V any](name, dir string) (out *Skiplist[K, V], err error) {
	var s Skiplist[K, V]
	s.fullPath = path.Join(dir, name)
	if err = os.MkdirAll(s.fullPath, 0744); err != nil {
		return
	}

	if s.floor, err = newFloor[K, V](s.fullPath); err != nil {
		return
	}

	if s.levels, err = openLevels[K](s.fullPath); err != nil {
		return
	}

	s.incrementEvery = 4
	out = &s
	return
}

type Skiplist[K Key, V any] struct {
	// Levels, starting at the top level
	levels []*level[K]
	floor  *floor[K, V]

	incrementEvery int
	fullPath       string
}

func (s *Skiplist[K, V]) Get(key K) (value V, err error) {
	seekIndex := s.getSeekIndex(key)
	return s.getMatch(seekIndex, key)
}

func (s *Skiplist[K, V]) Insert(key K, val V) (err error) {
	e := makeEntry(key, val)
	seekIndex := s.getSeekIndex(key)
	index, match := s.floor.getIndex(seekIndex, key)
	if match {
		return s.floor.Set(index, e)
	}

	if err = s.floor.InsertAt(index, e); err != nil {
		return
	}

	if s.floor.count%s.incrementEvery != 0 {
		return
	}

	lastIndex := index

	var nextLevel int
	for {
		var l *level[K]
		if l, err = s.getLevel(nextLevel); err != nil {
			return
		}

		if lastIndex, err = l.Insert(lastIndex, makeEntry(key, lastIndex)); err != nil {
			return
		}

		if l.count%s.incrementEvery != 0 {
			return
		}

		nextLevel++
	}
}

func (s *Skiplist[K, V]) getLevel(n int) (l *level[K], err error) {
	if len(s.levels) > n {
		return s.levels[n], nil
	}

	if l, err = newLevel[K](s.fullPath, n); err != nil {
		return
	}

	s.levels = append(s.levels, l)
	return
}

func (s *Skiplist[K, V]) getMatch(seekIndex int, key K) (value V, err error) {
	var match bool
	if value, match = s.floor.GetMatch(seekIndex, key); !match {
		err = ErrNotFound
		return
	}

	return
}

func (s *Skiplist[K, V]) getSeekIndex(key K) (seekIndex int) {
	for _, l := range s.levels {
		seekIndex = l.GetSeekIndex(seekIndex, key)
	}

	return
}

func (s *Skiplist[K, V]) printTree() {
	var (
		i int
		l *level[K]
	)
	for i := len(s.levels) - 1; i > -1; i-- {
		l = s.levels[i]
		l.printLayer()
	}

	if i > 0 {
		fmt.Print("\n")
	}

	s.floor.printLayer()
}
