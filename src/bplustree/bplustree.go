package bplustree

import (
  "unsafe"
  //"sort"
  "fmt"
  "bytes"
  "encoding/binary"
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
  Ptr      uint32      /* page number or Offset of the page start of a payload */
  Key      uint32      /* The key for Payload*/
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

  // get parent page
  pgHdr := bpTree.pPager.Read(pg.parent())
  if pgHdr == nil {
    panic("")
  }
  ppg := (*MemPage)(unsafe.Pointer(pgHdr.GetHeader()))

  for {
    ok, newpg = ppg.insert(bpTree, &Cell{Key: newpg.maxkey, Ptr: newpg.pgno})
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
  plbuf := new(bytes.Buffer)
  clbuf := new(bytes.Buffer)
  pgHdr := bpTree.pPager.Fetch(page.pgno)
  pagebuf := *(*[]byte)(unsafe.Pointer(pgHdr.pBulk))
  plof := uint16(unsafe.Sizeof(*&MemPage{})) + page.nfree +
          uint16(unsafe.Sizeof(*&Cell{})) * page.ncell
  clof := uint16(unsafe.Sizeof(*&MemPage{})) +
          uint16(unsafe.Sizeof(*&Cell{})) * page.ncell

  switch data.(type){
    case *Cell:
      if page.flag == INTERPAGE &&
         page.nfree >= uint16(unsafe.Sizeof(*&Cell{})) {

         binary.Write(clbuf, binary.LittleEndian, data.(*Cell).Ptr)
         binary.Write(clbuf, binary.LittleEndian, data.(*Cell).Key)
         copy(pagebuf[clof:], clbuf.Bytes())
         // write page head to page data
         page.nfree -= uint16(unsafe.Sizeof(*&Cell{}))
         page.ncell += 1
         if page.maxkey < data.(*Cell).Key {
           page.maxkey = data.(*Cell).Key
         }
         pgHdr.WriteHeader(page.flag, page.ncell, page.nfree, page.pgno, page.ppgno, page.maxkey)
         fmt.Printf("free=%d len=%d cap=%d slice=%v\n",clof, len(pagebuf),cap(pagebuf),pagebuf[0:1024])

         childHdr := bpTree.pPager.Fetch(data.(*Cell).Ptr)
         child := (*MemPage)(unsafe.Pointer(childHdr.GetHeader()))
         if page.pgno != child.ppgno {
           child.ppgno = page.pgno
           childHdr.WriteHeader(child.flag, child.ncell, child.nfree, child.pgno, child.ppgno, child.maxkey)
         }
         return true, nil
      } else {
        //key, newpg :=split(pg)
        newpg := bpTree.NewPage()
        newHdr := bpTree.pPager.Fetch(page.pgno)
        newPgbuf := *(*[]byte)(unsafe.Pointer(pgHdr.pBulk))

        of1 := uint16(unsafe.Sizeof(*&MemPage{}))
        of2 := of1 + uint16(unsafe.Sizeof(*&Cell{})) * (page.ncell/2)
        copy(newPgbuf[of1:], pagebuf[of2:newHdr.pCache.szPage])
        //update page info
        newpg.maxkey = page.maxkey
        newpg.flag = INTERPAGE
        newpg.ppgno = page.ppgno
        newpg.ncell = page.ncell/2
        newpg.nfree = newHdr.pCache.szPage - uint16(unsafe.Sizeof(*&PgHead{})) - uint16(unsafe.Sizeof(*&Cell{})) * newpg.ncell

        page.maxkey = page.cellptr(page.ncell/2).key
        page.ncell = page.ncell - newpg.ncell

        newHdr.WriteHeader(newpg.flag, newpg.ncell, newpg.nfree, newpg.pgno, newpg.ppgno, newpg.maxkey)
        pgHdr.WriteHeader(page.flag, page.ncell, page.nfree, page.pgno, page.ppgno, page.maxkey)

        fmt.Printf("free=%d len=%d cap=%d slice=%v\n",clof, len(newPgbuf),cap(newPgbuf),newPgbuf[0:1024])
        return false, newpg
      }
    case *Payload:
      if page.flag == LEAFPAGE  &&
         page.nfree >= (data.(*Payload).Size + uint16(unsafe.Sizeof(*&Cell{}))) {

         // write payload to page date
        binary.Write(plbuf, binary.LittleEndian, data.(*Payload).Key)
        binary.Write(plbuf, binary.LittleEndian, data.(*Payload).Size)
        for _, v := range data.(*Payload).Entrys {
            err := binary.Write(plbuf, binary.LittleEndian, v)
            if err != nil {
                fmt.Println("binary.Write failed:", err)
            }
        }
        fmt.Println("after write ，buf is:",plbuf.Bytes())
        copy(pagebuf[plof-data.(*Payload).Size:], plbuf.Bytes())
        // write cell for this payload to page date
        binary.Write(clbuf, binary.LittleEndian, data.(*Payload).Key)
        binary.Write(clbuf, binary.LittleEndian, data.(*Payload).Size)
        copy(pagebuf[clof:], clbuf.Bytes())

        // write page head to page data
        page.nfree -= (data.(*Payload).Size + uint16(unsafe.Sizeof(*&Cell{})))
        page.ncell += 1
        if page.maxkey < data.(*Payload).Key {
          page.maxkey = data.(*Payload).Key
        }
        pgHdr.WriteHeader(page.flag, page.ncell, page.nfree, page.pgno, page.ppgno, page.maxkey)
        fmt.Printf("free=%d len=%d cap=%d slice=%v\n",clof, len(pagebuf),cap(pagebuf),pagebuf[0:1024])
        return true, nil
      } else {
        //key, newpg :=split(pg)
        newpg := bpTree.NewPage()
        //update page info
        newpg.maxkey = page.maxkey

        page.maxkey = page.cellptr(page.ncell/2).key
        page.ncell = page.ncell/2

        return false, newpg
      }
  }

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

  fmt.Printf("no such type\n")
  return true
}

func (page *MemPage) cellptr(i uint16) *Cell {
  return (*Cell)(unsafe.Pointer(uintptr(unsafe.Pointer(page)) +
  uintptr(unsafe.Sizeof(&MemPage{})) +
  uintptr(i)*unsafe.Sizeof(&Cell{})))
}
