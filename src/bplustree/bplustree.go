package bplustree

import (
  "unsafe"
)

const (
  LEAFPAGE = 0
  INTERPAGE = 1
)
/*
** An instance of this object represents a single database file.
**
** A single database file can be in use at the same time by two
** or more database connections.  When two or more connections are
** sharing the same database file, each connection has it own
** private Btree object for the file and each of those Btrees points
** to this one BPlusTree object.
*/
type BPlusTree struct {
  pPager *Pager           /* The page cache */
  page *PgHead          /* First page of the database */
}

type MemPage PgHead

/* The basic idea is that each page of the file contains N database
** entries and N+1 pointers to subpages.
**
**   --------------------------------------------------------------
**   |  Ptr(0) | Key(0) | Ptr(1) | Key(1) | ... | Key(N) | Ptr(N) |
**   --------------------------------------------------------------
*/
type Cell struct {
  ptr      uint32      /* page number or Offset of the page start of a payload */
  key      uint32      /* The key for Payload*/
}

/* DocId1 DocId2 ...
**   -----------------------------------------------------------------
**   |  key | DocSiz | DocId1 | DocId3 | ... | DocId(N-1) | DocId(N) |
**   -----------------------------------------------------------------
 */
type Payload struct {
  key     uint32             /* value in the unpacked key */
  size    uint16             /* Number of values.  Might be zero */
  entrys  *[]byte            /* fot data compress */
}

func (bpTree *BPlusTree) Open(dbName string, dbSize int) {
  //bpTree.pPager = &Pager{}
  pPager := bpTree.pPager
  pPager.Open(dbName, dbSize)

  bpTree.page = (*PgHead)(unsafe.Pointer(pPager.Fetch(0)))
}

func (bpTree *BPlusTree) Insert(pl *Payload) {
  of, pg := bpTree.Search(pl.key)
  if pg == nil {
    return
  }

  ok, key, newpg := pg.insert(of, pl)
  if ok {
    return
  }

  // get parent page
  pgHdr := bpTree.pPager.Read(pg.parent())
  if pgHdr == nil {
    panic("")
  }
  ppg := (*MemPage)(unsafe.Pointer(pgHdr.GetPageHeader()))

  for {
    ok, key, newpg = ppg.insert(0, &Cell{key: key,ptr: newpg.pgno})
    if ok {
      return
    }

    if ppg.pgno == bpTree.page.pgno {
      // alloc new root page for bplustree and update bplustree page
      rootpage := &MemPage{}
      bpTree.page = (*PgHead)(unsafe.Pointer(rootpage))
      // insert new page cell
      rootpage.insert(0, &Cell{key: key,ptr: newpg.pgno})

      // insert origin page cell
      rootpage.insert(0, &Cell{key: key,ptr: ppg.pgno})
      return
    }

    // get parent page
    pgHdr = bpTree.pPager.Read(ppg.parent())
    if pgHdr == nil {
      panic("")
    }
    ppg = (*MemPage)(unsafe.Pointer(pgHdr.GetPageHeader()))
  }
}

func (bpTree *BPlusTree) Search(key uint32) (uint16, *MemPage) {
  curr := (*MemPage)(unsafe.Pointer(bpTree.page))
  for {
    switch curr.flag {
    case LEAFPAGE:
      of, ok := find(curr, key)
      if !ok {
        return 0, curr
      }
      return of, curr
    case INTERPAGE:
      pgno, _ := find(key)
      // curr = bpTree.hm[pgno]
      // pager should load page and cached
      pgHdr := bpTree.pPager.Read(pgno)
      if pgHdr == nil {
        panic("")
      }
      curr = (*MemPage)(unsafe.Pointer(pgHdr.GetPageHeader()))
    default:
      panic("no such flag!")
    }
  }
}

func (page *MemPage) insert(offset uint16, data interface{}) (bool, uint32, *MemPage){
  ok := page.full(data)
  if !ok {
    return true, 0, nil
  }

  //key, newpg :=split(pg)
  newpg := newpage()
  //update page info
  newpg.maxkey = page.maxkey
  pcell := page.cellptr()

  page.maxkey = pcell[page.ncell/2].key
  page.ncell = page.ncell/2

  return false, key, newpg
}

func (page *MemPage) find(key int) (int, bool) {
  if  page.pgno == 0 && page.ncell == 0 {
    return 1, false
  }

  pcell := page.cellptr()

  cmp := func (i int) bool {
    return pcell[i].key >= key
  }

  i := sort.Search(page.ncell, cmp)

  if page.flag == INTERPAGE {
    return pcell[i].ptr, true
  }

  if i <= page.ncell && pcell[i].key == key {
    return pcell[i].ptr, true
  }

  return nil, false
}

func newpage() *MemPage{
  page := &MemPage{}
  // alloc page in cache
  // page.cell = cache.allocpage()[100]
}

func (page *MemPage) parent() uint32 {
  return page.ppgno
}

func (page *MemPage) setparent(uint32 pgno) {
  page.ppgno = pgno
}

func (page *MemPage) full(data interface{}) bool {
  // insert payload sorted by pl.key todo later...
  switch data.(type){
  case *Cell:
    if page.flag == INTERPAGE {
      return page.free < size(Cell)
    }
    panic("full error")
  case *PlayLoad:
    if page.flag == LEAFPAGE {
      return page.free < (pl.size + unsafe.Sizeof(&Cell{}))
    }
    panic("full error")
  }
}

func (page *MemPage) cellptr() uint32 {
  return (*Cell)(unsafe.Pointer(uintptr(unsafe.Pointer(page)) + uintptr(unsafe.Sizeof(&MemPage{}))))
}
