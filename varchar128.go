package bssl

import "slices"

func MakeVarchar128(str string) (out Varchar128) {
	copy(out[:], str)
	return
}

type Varchar128 [128]byte

func (v Varchar128) String() string {
	return toString(v[:])
}

func (v Varchar128) Compare(in Varchar128) (result int) {
	return slices.Compare(v[:], in[:])
}
