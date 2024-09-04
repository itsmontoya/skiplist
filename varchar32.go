package skiplist

var _ Key = Varchar32{}

func MakeVarchar32(str string) (out Varchar32) {
	copy(out[:], str)
	return
}

type Varchar32 [32]byte

func (v Varchar32) String() string {
	return toString(v[:])
}

func (v Varchar32) Compare(in any) (result int) {
	b := *(in.(*Varchar32))
	for i := 0; i < 32; i++ {
		aV := v[i]
		bV := b[i]
		switch {
		case aV < bV:
			return -1
		case aV > bV:
			return 1
		}
	}

	return 0
}
