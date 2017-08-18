package bplustree

import (
  "unsafe"
  //"sort"
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
  page *MemPage          /* First page of the database */
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

func (bpTree *BPlusTree) Open(dbName string, dbSize uint16) {
  bpTree.pPager = &Pager{}
  bpTree.pPager.Create(dbName, dbSize)

  pgHdr := bpTree.pPager.Fetch(0)
  bpTree.page = (*MemPage)(unsafe.Pointer(pgHdr.GetPageHeader()))
  println("bplustree Open", len(*(*[]byte)(unsafe.Pointer(pgHdr.pBulk))))
}

func (bpTree *BPlusTree) Insert(pl *Payload) {
  _, pg := bpTree.Search(pl.key)
  if pg == nil {
    return
  }

  ok, newpg := pg.insert(bpTree, pl)
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
    ok, newpg = ppg.insert(bpTree, &Cell{key: newpg.maxkey,ptr: newpg.pgno})
    if ok {
      return
    }

    if ppg.pgno == bpTree.page.pgno {
      pg := bpTree.NewPage()
      // move rootpage data to newpage

      // insert new page cell
      bpTree.page.insert(bpTree, &Cell{key: newpg.maxkey,ptr: newpg.pgno})

      // insert origin page cell
      bpTree.page.insert(bpTree, &Cell{key: pg.maxkey,ptr: pg.pgno})
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

func (bpTree *BPlusTree) Search(key uint32) (uint32, *MemPage) {
  curr := (*MemPage)(unsafe.Pointer(bpTree.page))
  for {
    switch curr.flag {
    case LEAFPAGE:
      of, ok := curr.find(key)
      if !ok {
        return 0, curr
      }
      return of, curr
    case INTERPAGE:
      pgno, _ := curr.find(key)
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

func (page *MemPage) insert(bpTree *BPlusTree, data interface{}) (bool, *MemPage){
  ok := page.full(data)
  if !ok {
    return true, nil
  }

  //key, newpg :=split(pg)
  newpg := bpTree.NewPage()
  //update page info
  newpg.maxkey = page.maxkey

  page.maxkey = page.cellptr(page.ncell/2).key
  page.ncell = page.ncell/2

  return false, newpg
}

func (page *MemPage) find(key uint32) (uint32, bool) {
  /*if  page.pgno == 0 && page.ncell == 0 {
    return 1, false
  }

  pcell := page.cellptr(0)

  cmp := func (i int) bool {
    return page.cellptr(i).key >= key
  }

  i := sort.Search(page.ncell, cmp)

  if page.flag == INTERPAGE {
    return page.cellptr(i).ptr, true
  }

  if i <= page.ncell && page.cellptr(i).key == key {
    return page.cellptr(i).ptr, true
  }*/

  return 0, false
}

func (bpTree *BPlusTree) NewPage() *MemPage{
  pPager := bpTree.pPager
  pPager.numPage += 1
  pgHdr := pPager.Fetch(pPager.numPage)
  return (*MemPage)(unsafe.Pointer(pgHdr.GetPageHeader()))
}

func (page *MemPage) parent() uint32 {
  return page.ppgno
}

func (page *MemPage) setparent(pgno uint32) {
  page.ppgno = pgno
}

func (page *MemPage) full(data interface{}) bool {
  // insert payload sorted by pl.key todo later...
  switch data.(type){
  case *Cell:
    if page.flag == INTERPAGE {
      return page.nfree < uint16(unsafe.Sizeof(&Cell{}))
    }
    panic("full error")
  case *Payload:
    if page.flag == LEAFPAGE {
      return page.nfree < (data.(Payload).size + uint16(unsafe.Sizeof(&Cell{})))
    }
    panic("full error")
  }
  panic("tyoe error")
  return true
}

func (page *MemPage) cellptr(i uint16) *Cell {
  return (*Cell)(unsafe.Pointer(uintptr(unsafe.Pointer(page)) +
  uintptr(unsafe.Sizeof(&MemPage{})) +
  uintptr(i)*unsafe.Sizeof(&Cell{})))
}
