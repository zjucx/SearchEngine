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

func search(n node, key int) (*kv, int, *leafNode) {
  curr := bt.root
  index := -1
  for {
    switch t := curr.(type) {
    case *leafNode:
      i, ok := t.find(key)
      if !ok {
        return "", index, t
      }
      return t.kvs[i].value, index, t
    case *interNode:
      i, _ = t.find(key)
      index = i
      curr = t.kcs[i].child
    default:
      panic("")
    }
  }
}
