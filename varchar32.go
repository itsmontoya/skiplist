package bssl

import "slices"

func MakeVarchar32(str string) (out Varchar32) {
	copy(out[:], str)
	return
}

type Varchar32 [32]byte

func (v Varchar32) String() string {
	return toString(v[:])
}

func (v Varchar32) Compare(in Varchar32) (result int) {
	return slices.Compare(v[:], in[:])
}
