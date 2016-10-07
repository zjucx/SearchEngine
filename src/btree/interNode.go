package bplustree

type kc struct {
  key string
  child node
}

type interNode struct {
  kcs [MaxKC]kc
  count int
  parent *interNode
}

func newInterNode(p *interNode, largestChild node) *interNode {
  return &interNode{
    parent : p,
    count : 1,
    kcs[0].child : largestChild != nil ? largestChild : nil,
  }
}

func (in *interNode) find(key int) (int, bool) {
  cmp := func (i int) bool {
    return in.kcs[i].key >= key
  }

  i := sort.Search(in.count, cmp)

  return i, true
}
