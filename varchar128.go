package skiplist

import "slices"

var _ Key = Varchar128{}

func MakeVarchar128(str string) (out Varchar128) {
	copy(out[:], str)
	return
}

type Varchar128 [128]byte

func (v Varchar128) String() string {
	return toString(v[:])
}

func (v Varchar128) Compare(in any) (result int) {
	b := *(in.(*Varchar128))
	return slices.Compare(v[:], b[:])
}
