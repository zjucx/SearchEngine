package bplustree

import (
  "unsafe"
  //"sort"
  "fmt"
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
  Key     uint32             /* value in the unpacked key */
  Size    uint16             /* Number of values.  Might be zero */
  Entrys  []uint32            /* fot data compress */
}

func (bpTree *BPlusTree) Open(dbName string, dbSize uint16) {
  bpTree.pPager = &Pager{}
  bpTree.pPager.Create(dbName, dbSize)

  pgHdr := bpTree.pPager.Fetch(0)
  bpTree.page = (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))
  println("bplustree Open", len(*(*[]byte)(unsafe.Pointer(pgHdr.pBulk))))
}

func (bpTree *BPlusTree) Insert(pl *Payload) {
  pg := bpTree.Search(pl.Key)
  if pg == nil {
    return
  }
  fmt.Printf("flag=%d pgno=%d\n",pg.flag, pg.pgno)

  ok, newpg := pg.insert(bpTree, pl)
  if ok {
    return
  }
fmt.Printf("flag=%d pgno=%d\n",pg.flag, pg.pgno)
  // get parent page
  pgHdr := bpTree.pPager.Read(pg.parent())
  if pgHdr == nil {
    panic("")
  }
  ppg := (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))

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
    ppg = (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))
  }
}

func (bpTree *BPlusTree) Search(key uint32) (*MemPage) {
  curr := (*MemPage)(unsafe.Pointer(bpTree.page))

  for {
    switch curr.flag {
    case LEAFPAGE:
      return curr
    case INTERPAGE:
      pgno, _ := curr.find(key)
      // curr = bpTree.hm[pgno]
      // pager should load page and cached
      pgHdr := bpTree.pPager.Read(pgno)
      if pgHdr == nil {
        panic("")
      }
      curr = (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))
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
  switch page.flag {
    case LEAFPAGE:
      return 0, true
    case INTERPAGE:
      if page.ncell == 0 {
        return 1, false
      }

      bIdx := uint16(0)
      eIdx := page.ncell
      for bIdx != eIdx {
        if key > page.cellptr(bIdx).key {
          bIdx = (bIdx + eIdx)/2
        } else {
          eIdx = (bIdx + eIdx)/2
        }
      }
      return page.cellptr(bIdx).ptr, true
    default:
      panic("no such flag!")
  }
}

func (bpTree *BPlusTree) NewPage() *MemPage{
  pPager := bpTree.pPager
  pPager.numPage += 1
  pgHdr := pPager.Fetch(pPager.numPage)
  return (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))
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
        return page.nfree < uint16(unsafe.Sizeof(*&Cell{}))
      }
    case *Payload:
      if page.flag == LEAFPAGE {
        return page.nfree < (data.(*Payload).Size + uint16(unsafe.Sizeof(*&Cell{})))
      }
  }
  fmt.Printf("no such type\n")
  return true
}

func (page *MemPage) cellptr(i uint16) *Cell {
  return (*Cell)(unsafe.Pointer(uintptr(unsafe.Pointer(page)) +
  uintptr(unsafe.Sizeof(&MemPage{})) +
  uintptr(i)*unsafe.Sizeof(&Cell{})))
}
