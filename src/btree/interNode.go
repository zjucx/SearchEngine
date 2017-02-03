package bplustree

type kc struct {
  key string
  child node
}

type InterNode struct {
  kcs [MaxKC]kc
  count int
  parent *InterNode
}

func NewInterNode(p *InterNode, largestChild node) *InterNode {
  return &InterNode{
    parent : p,
    count : 1,
    kcs[0].child : largestChild != nil ? largestChild : nil,
  }
}

func (in *InterNode) Insert(key int, value node) (*InterNode, int, bool) {
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

func split(in  *InterNode) (*InterNode, int) {
  midIndex := MaxKC/2
  newInter := &InterNode{
    count : midIndex,
  }

  copy(newInter.kcs[0:], in.kcs[midIndex:MaxKC])

  newInter.count = MaxKC - midIndex

  return newInter, newInter.kcs[0].key
}

func (in *InterNode) find(key int) (int, bool) {
  cmp := func (i int) bool {
    return in.kcs[i].key >= key
  }

  i := sort.Search(in.count, cmp)

  return i, true
}

func (in *InterNode) parent() *InterNode {return in.parent}
func (in *InterNode) setParent(in *InterNode) {in.parent = in}
func (in *InterNode) full() bool {return in.count == MaxKC}
