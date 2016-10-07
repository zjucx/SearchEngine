package bplustree

const (
  MaxKV = 2
  MaxKC = 2
)

type node interface {
  find(key int) (int, bool)
  parent() *interNode
  setParent(*interNode)
  full() bool
}
