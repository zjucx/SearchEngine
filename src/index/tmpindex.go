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

type item struct {
  wordid    int
  docid     int
}

/*
 * the tmp index file(.tmp) 's data orgnized by
 */
 type tmpIndexBuf struct {
   buf [maxbufsize]byte		/* 输入缓冲区 */
   itemptr *item
   length int		/* 缓冲区当前有多少个数 */
   offset int	/* 缓冲区读到了文件的哪个位置 */
   idx int		/* 缓冲区的指针 */
 }


 /*
  * the mian enter of algorithm of quicksort
  */
func (idx *tmpIndex)quickSort(s, t int) {
  m := idx.split()
  idx.quickSort(idx, s, m-1)
  idx.quickSort(idx, m+1, t)
}
/*
 * the split part of algorithm of quicksort
 */
func (idx *tmpIndex)split(idx []int, s, t int) int {
  for i, j:= s; i < t; i++ {
    if idx.Less(i, t) {
      idx.Swap(i, j)
      j++
    }
  }
  idx.Swap(j, t)
  return j
}
func (idx *tmpIndex)Swap(i, j int) {
  idx.tmpidxptr[i], idx.tmpidxptr[j] = idx.tmpidxptr[j], idx.tmpidxptr[i]
}
func (idx *tmpIndex) Less(i, j int) bool {
  if idx.tmpidxptr[i].wordid < idx.tmpidxptr[j].wordid {
    return true
  }

  if idx.tmpidxptr[i].docid < idx.tmpidxptr[j].docid {
    return true
  }

  return false
}
