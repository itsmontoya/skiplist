package skiplist

import "slices"

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
	b := in.(Varchar32)
	return slices.Compare(v[:], b[:])
}
