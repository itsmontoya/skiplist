package skiplist

type Key interface {
	Compare(any) int
}
