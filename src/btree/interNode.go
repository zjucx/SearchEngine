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

func (in *interNode) Insert(key int, value string) (*interNode, int, bool) {
  index, _ := find(key)
  if !in.full() {
    copy(in.kcs[index:], in.kcs[index+1 : count])
    in.count++
    return nil, key, false
  }

  
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
