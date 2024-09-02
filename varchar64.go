package bssl

import "slices"

func MakeVarchar64(str string) (out Varchar64) {
	copy(out[:], str)
	return
}

type Varchar64 [64]byte

func (v Varchar64) String() string {
	return toString(v[:])
}

func (v Varchar64) Compare(in Varchar64) (result int) {
	return slices.Compare(v[:], in[:])
}
