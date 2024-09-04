package skiplist

import "slices"

var _ Key = Varchar16{}

func MakeVarchar16(str string) (out Varchar16) {
	copy(out[:], str)
	return
}

type Varchar16 [16]byte

func (v Varchar16) String() string {
	return toString(v[:])
}

func (v Varchar16) Compare(in any) (result int) {
	b := *(in.(*Varchar16))
	return slices.Compare(v[:], b[:])
}
