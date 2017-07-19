package bplustree

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
