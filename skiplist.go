package bssl

import (
	"errors"
	"os"
	"path"
)

var (
	ErrNotFound = errors.New("entry not found")
)

func New[K Key[any], V any](name, dir string) (out *Skiplist[K, V], err error) {
	var s Skiplist[K, V]
	fullPath := path.Join(dir, name)
	if err = os.MkdirAll(fullPath, 0744); err != nil {
		return
	}

	if s.floor, err = newFloor[K, V](fullPath); err != nil {
		return
	}

	if s.levels, err = openLevels[K](fullPath); err != nil {
		return
	}

	out = &s
	return
}

type Skiplist[K Key[any], V any] struct {
	// Levels, starting at the top level
	levels []*level[K]
	floor  *floor[K, V]
}

func (s *Skiplist[K, V]) Get(key K) (value V, err error) {
	seekIndex := s.getSeekIndex(key)
	return s.getMatch(seekIndex, key)
}

func (s *Skiplist[K, V]) getMatch(seekIndex int, key K) (value V, err error) {
	var match bool
	if value, match = s.floor.getMatch(seekIndex, key); !match {
		err = ErrNotFound
		return
	}

	return
}

func (s *Skiplist[K, V]) getSeekIndex(key K) (seekIndex int) {
	for _, l := range s.levels {
		seekIndex = l.getIndex(seekIndex, key)
	}

	return
}
