package skiplist

import "slices"

var _ Key = Varchar64{}

func MakeVarchar64(str string) (out Varchar64) {
	copy(out[:], str)
	return
}

type Varchar64 [64]byte

func (v Varchar64) String() string {
	return toString(v[:])
}

func (v Varchar64) Compare(in any) (result int) {
	b := in.(Varchar64)
	return slices.Compare(v[:], b[:])
}
