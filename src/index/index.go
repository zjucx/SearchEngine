package index

import (
  "fmt"
  "unsafe"
  "container/list"
)

const (
	tindex     uint8 = 0
	tdata      uint8 = 1
)

/*
 * the final index file(.index) and the data file (.data)'s data
 * orgnized by page
 */
type page struct {
  pgid uint32
  nextid uint32
  preid uint32
  pgtype uint8       //dacid data page and index page
  curSize uint32
  data uintptr      //save data for docid or index page
}

/*
 * the tmp index file(.tmp) 's data orgnized by
 */
 type tmpindex struct {
   dictid    int
   docid     int
 }

//index
type index struct {
  tmpidxptr *tmpindex
  pg *page
}

func (idx *index)quickSort(idx []int, s, t int) {
  m := split()
  quickSort(idx, s, m-1)
  quickSort(idx, m+1, t)
}

func (idx *index)split(idx []int, s, t int) int {
  for i, j:= s; i < t; i++ {
    if idx[i] < idx[t] {
      idx.swap(i, j)
      j++
    }
  }
  swap(idx+j, idx+t)
  return j
}
func (idx *index)swap(i, j int) {
  idx.tmpidxptr[i], idx.tmpidxptr[j] = idx.tmpidxptr[j], idx.tmpidxptr[i]
}
