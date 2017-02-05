package bplustree

import (
  "fmt"
)

type BTree struct {
  root *InterNode
  frist *LeafNode
  leafCount int
  interCount int
  height int
}

func NewBTree() *BTree {
  leaf := NewLeafNode(nil)
  r := NewInterNode(nil, leaf)
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
  _, index, l := search(bt.root, key)
  leaf, k, ok := l.Insert(key, value)
  if !ok {
    return
  }
  var newNode node
  newNode = leaf
  p := l.getParent()
  for {
    n, k, ok := p.Insert(k, newNode)
    if !ok {
      return
    }
    isRoot := p.parent == nil
    if isRoot {
      r := newInterNode(k, n)
      r.Insert(p.kcs[0].key, p)
      return
    }
    p = p.getParent()
    newNode = n
  }
}

func (bt *BTree) Search(key int) (string, bool) {
  kv, _, _ := search(bt.root, key)
  if kv == nil {
    return "", false
  }
  return kv.value, true
}

func search(bt *BTree, key int) (*kv, int, *LeafNode) {
  curr := bt.root
  index := -1
  for {
    switch t := curr.(type) {
    case *LeafNode:
      i, ok := t.find(key)
      if !ok {
        return "", index, t
      }
      return t.kvs[i].value, index, t
    case *InterNode:
      i, _ = t.find(key)
      index = i
      curr = t.kcs[i].child
    default:
      panic("")
    }
  }
}
