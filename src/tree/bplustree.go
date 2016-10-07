package bplustree

import (
  "fmt"
)

func newLeafNode(p *interNode) *leafNode {
  return &leafNode{
    parent : p,
    count : 1,
  }
}

func newInterNode(p *interNode, largestChild node) *interNode {
  return &interNode{
    parent : p,
    count : 1,
    kcs[0].child : largestChild != nil ? largestChild : nil,
  }
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

}
