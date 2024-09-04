package skiplist

import (
	"errors"
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func openLevels[K Key](fullPath string) (out levels[K], err error) {
	if err = walkLevels(fullPath, func(filepath string, parsed int) (err error) {
		var l *level[K]
		if l, err = newLevel[K](filepath, parsed); err != nil {
			return
		}

		out = append(out, l)
		return
	}); err != nil {
		return
	}

	sort.Slice(out, func(a, b int) (more bool) {
		lA := out[a]
		lB := out[b]
		return lA.level > lB.level
	})

	return
}

func walkLevels(fullPath string, fn func(filepath string, parsed int) error) (err error) {
	if err = filepath.Walk(fullPath, func(filepath string, info fs.FileInfo, ierr error) (err error) {
		if ierr != nil {
			return
		}

		if info.IsDir() {
			return
		}

		base := path.Base(filepath)
		ext := path.Ext(base)
		if ext != ".level" {
			return
		}

		withoutExtension := strings.Replace(base, ext, "", 1)

		var parsed int
		if parsed, err = strconv.Atoi(withoutExtension); err != nil {
			return
		}

		return fn(filepath, parsed)
	}); err != nil {
		return
	}

	return
}

type levels[K Key] []*level[K]

func (ls levels[K]) Close() (err error) {
	var errs []error
	for _, l := range ls {
		if err := l.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
func (ls levels[K]) printLayers() {
	ls.iterateFromTopLevel(func(_ int, l *level[K]) bool {
		l.printLayer()
		return false
	})
}

func (ls levels[K]) iterateFromTopLevel(fn func(i int, l *level[K]) (end bool)) (ended bool) {
	for i := len(ls) - 1; i > -1; i-- {
		l := ls[i]
		if ended = fn(i, l); ended {
			return
		}
	}

	return false
}
