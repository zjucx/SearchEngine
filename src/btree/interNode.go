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

func (in *interNode) Insert(key int, value node) (*interNode, int, bool) {
  index, _ := find(key)
  if !in.full() {
    copy(in.kcs[index+1:], in.kcs[index : count])
    in.kcs[index].key = key
    in.kcs[index].child = value
    in.count++
    return nil, key, false
  }

  newInter, k := split(in)

  if key > k {
    newInter.Insert(key, value)
  } else {
    in.Insert(key, value)
  }
  return newInter, k, true
}

func split(in  *interNode) (*interNode, int) {
  midIndex := MaxKC/2
  newInter := &interNode{
    count : midIndex,
  }

  copy(newInter.kcs[0:], in.kcs[midIndex:MaxKC])

  newInter.count = MaxKC - midIndex

  return newInter, newInter.kcs[0].key
}

func (in *interNode) find(key int) (int, bool) {
  cmp := func (i int) bool {
    return in.kcs[i].key >= key
  }

  i := sort.Search(in.count, cmp)

  return i, true
}

func (in *interNode) parent() *interNode {return in.parent}
func (in *interNode) setParent(in *interNode) {in.parent = in}
func (in *interNode) full() bool {return in.count == MaxKC}
