package skiplist

import (
	"errors"
	"os"
	"path"
)

var (
	ErrNotFound = errors.New("entry not found")
)

func New[K Key, V any](name, dir string, size int64) (out *Skiplist[K, V], err error) {
	var s Skiplist[K, V]
	s.fullPath = path.Join(dir, name)
	if err = os.MkdirAll(s.fullPath, 0744); err != nil {
		return
	}

	if s.floor, err = newFloor[K, V](s.fullPath, size); err != nil {
		return
	}

	if s.levels, err = openLevels[K](s.fullPath); err != nil {
		return
	}

	s.incrementEvery = 40
	out = &s
	return
}

type Skiplist[K Key, V any] struct {
	// Levels, starting at the top level
	levels levels[K]
	floor  *floor[K, V]

	incrementEvery int
	fullPath       string
}

func (s *Skiplist[K, V]) Get(key K) (value V, err error) {
	seekIndexes := s.getSeekIndexes(&key)
	return s.getMatch(seekIndexes.Last(), &key)
}

func (s *Skiplist[K, V]) Set(key K, val V) (err error) {
	return s.set(key, val, true)
}

func (s *Skiplist[K, V]) SetNX(key K, val V) (err error) {
	return s.set(key, val, false)
}

func (s *Skiplist[K, V]) Close() (err error) {
	var errs []error
	if err = s.levels.Close(); err != nil {
		errs = append(errs, err)
	}

	if err = s.floor.Close(); err != nil {
		errs = append(errs, err)
	}

	return errors.Join(errs...)
}

func (s *Skiplist[K, V]) set(key K, val V, allowUpdate bool) (err error) {
	e := makeEntry(key, val)
	seekIndexes := s.getSeekIndexes(&key)
	index, match := s.floor.getIndex(seekIndexes.Last(), &key)
	switch {
	case !match:
		// No match, continue on
	case allowUpdate:
		// Match and allow update, set value at index
		return s.floor.Set(index, e)
	default:
		// Match and update not allowed, return
		return
	}

	if err = s.floor.InsertAt(index, e); err != nil {
		return
	}

	var topLevel int
	if s.floor.count%s.incrementEvery == 0 {
		if topLevel, err = s.insertReferences(key, index); err != nil {
			return
		}
	}

	return s.updateReferences(e.Key, topLevel, seekIndexes)
}

func (s *Skiplist[K, V]) updateReferences(key K, topLevel int, seekIndexes indexes) (err error) {
	s.levels.iterateFromTopLevel(func(i int, l *level[K]) (end bool) {
		if i > topLevel {
			return
		}

		l.IterateAfter(seekIndexes.At(i), key, func(index int, e Entry[K, int]) {
			e.Value++
			l.Set(index, e)
		})

		return false
	})

	return
}

func (s *Skiplist[K, V]) insertReferences(key K, lastIndex int) (nextLevel int, err error) {
	for {
		var l *level[K]
		if l, err = s.getLevel(nextLevel); err != nil {
			return
		}

		if lastIndex, err = l.Insert(0, makeEntry(key, lastIndex)); err != nil {
			return
		}

		nextLevel++

		if l.count%s.incrementEvery != 0 {
			return
		}
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

func (s *Skiplist[K, V]) getMatch(seekIndex int, key *K) (value V, err error) {
	var match bool
	if value, match = s.floor.GetMatch(seekIndex, key); !match {
		err = ErrNotFound
		return
	}

	return
}

func (s *Skiplist[K, V]) getSeekIndexes(key *K) (seekIndexes indexes) {
	var seekIndex int
	s.levels.iterateFromTopLevel(func(i int, l *level[K]) (end bool) {
		seekIndex = l.GetSeekIndex(seekIndex, key)
		seekIndexes = append(seekIndexes, seekIndex)
		return false
	})

	return
}

func (s *Skiplist[K, V]) printTree() {
	s.levels.printLayers()
	s.floor.printLayer()
}
