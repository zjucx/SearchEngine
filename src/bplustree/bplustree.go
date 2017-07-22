package bplustree

import (

)

const (
  LEAFPAGE = 0
  INTERPAGE = 1
)

func (bptree *BPlusTree) Insert(pl *PlayLoad) {

}

func (bptree *BPlusTree) search(pl *PlayLoad) (uint16, *MemPage) {
  curr := bptree.pPage
  for {
    switch t := curr.ph.flag {
    case LEAFPAGE:
      i, ok := t.find(curr, pl.key)
      if !ok {
        return "", index, curr
      }
      return t.kvs[i].value, index, t
    case INTERPAGE:
      i, _ = t.find(key)
      curr = bptree.hm[i]
    default:
      panic("")
    }
  }
}

func (bptree *BPlusTree) find(p *MemPage, key int) (int, bool) {
  cmp := func (i int) bool {
    return in.kcs[i].key >= key
  }

  i := sort.Search(in.count, cmp)

  return i, true
}

func (bptree *BPlusTree) parent(p *MemPage) *MemPage {
  return bptree.hm[p.ph.pgno]
}
func (bptree *BPlusTree) setparent(uint32 pgno) {
  p.ph.parent = pgno
}
func (bptree *BPlusTree) full(p *MemPage, pl *PlayLoad) bool {
  if p.ph.flag == 1 {
    return p.ph.nFree > (pl.size + size(Cell))
  }
}
