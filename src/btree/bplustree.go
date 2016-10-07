package bplustree

import (
  "fmt"
)

type BTree struct {
  root *interNode
  frist *leafNode
  leafCount int
  interCount int
  height int
}

func newBTree() *BTree {
  leaf := newLeafNode(nil)
  r := newInterNode(nil, leaf)
  leaf.parent = r

  return &BTree{
    root : r,
    frist : leaf,
    leafCount : 1,
    interCount : 1,
    height : 2,
  }
}


func (bt *BTree) Insert(key int, value string) {

}

func (bt *BTree) Search(key int) (string, bool) {
  curr := bt.root
  for {
    switch t := curr.(type) {
    case *leafNode:
    case *interNode:
      curr = t.
    }
  }
}
