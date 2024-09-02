package skiplist

import (
	"io/fs"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func openLevels[K Key](fullPath string) (out []*level[K], err error) {
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
