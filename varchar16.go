package bssl

import "slices"

func MakeVarchar16(str string) (out Varchar16) {
	copy(out[:], str)
	return
}

type Varchar16 [16]byte

func (v Varchar16) String() string {
	return toString(v[:])
}

func (v Varchar16) Compare(in Varchar16) (result int) {
	return slices.Compare(v[:], in[:])
}
