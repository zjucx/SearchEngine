package bplustree

import (
  "container/list"
  "sort"
)

type Key int

type kv struct {
  key Key
  value *list
}

type LeafNode struct {
  kvs [MaxKV]kv
  count int
  next *LeafNode
  parent *InterNode
}


func NewLeafNode(p *interNode) *LeafNode {
  return &LeafNode{
    parent : p,
    count : 1,
  }
}

func (l *LeafNode) find(key int) (int, bool) {
  cmp := func (i int) bool {
    return l.kvs[i].key >= key
  }

  i := sort.Search(l.count, cmp)

  if i < l.count && l.kvs[i].key == key {
    return i, true
  }

  return i, false
}

func (l *LeafNode) Insert(key int, value string) (*LeafNode, Key, bool) {
  index, ok := l.find(key)
  if !ok {
    l.kvs[index].value = value
    return nil, key, false
  }

  if !l.full() {
    copy(l.kvs[index+1:], l.kvs[index:l.count])
    l.kvs[index].key = key
    l.kvs[index].value = value
    l.count++
    return nil, key, false
  }

  newLeaf, k := split(l)

  if key > k {
    newLeaf.Insert(key, value)
  } else {
    l.Insert(key, value)
  }
  return newLeaf, k, true
}

func split(l *LeafNode) (*LeafNode, key) {
  midIndex := MaxKV/2
  newLeaf := &LeafNode{
    count : MaxKV - midIndex,
    next : l.next
  }

  copy(newLeaf.kvs[0:], l.kvs[midIndex:MaxKV])

  l.count = midIndex
  l.next = newLeaf

  return newLeaf, newLeaf.kvs[0].key
}

func (l *LeafNode) parent() *interNode {return l.parent}
func (l *LeafNode) setParent(in *interNode) {l.parent = in}
func (l *LeafNode) full() bool {return l.count == MaxKV}
