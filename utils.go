package skiplist

func findEnd(bs []byte) (i int) {
	var b byte
	for i, b = range bs {
		if b == 0 {
			return
		}
	}

	i++
	return
}

func toString(bs []byte) (str string) {
	i := findEnd(bs)
	return string(bs[:i])
}
