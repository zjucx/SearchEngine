package bplustree

import (
  "container/list"
  "sort"
)

type kv struct {
  key string
  value *list
}

type leafNode struct {
  kvs [MaxKV]kv
  count int
  next *leafNode
  parent *interNode
}


func newLeafNode(p *interNode) *leafNode {
  return &leafNode{
    parent : p,
    count : 1,
  }
}

func (l *leafnode) find(key int) (int, bool) {
  cmp := func (i int) bool {
    return l.kvs[i].key >= key
  }

  i := sort.Search(l.count, cmp)

  if i < l.count && l.kvs[i].key == key {
    return i, true
  }

  return i, false
}
