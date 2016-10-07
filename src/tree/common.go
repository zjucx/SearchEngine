package bplustree

import (
  "container/list"
)

const (
  MaxKV = 2
  MaxKC = 2
)

type node interface {

}

type kv struct {
  key string
  value *list
}

type kc struct {
  key string
  child node
}

type leafNode struct {
  kvs [MaxKV]kv
  count int
  next *leafNode
  parent *interNode
}

type interNode struct {
  kcs [MaxKC]kc
  count int
  parent *interNode
}

type BTree struct {
  root *interNode
  frist *leafNode
  leafCount int
  interCount int
  height int
}
