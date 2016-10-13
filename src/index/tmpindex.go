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

const K = 64
const maxbufsize = 4096

type item struct {
  wordid    int
  docid     int
}

/*
 * the tmp index file(.tmp) 's data orgnized by
 */
 type tmpIndexBuf struct {
   buf [maxbufsize]byte		/* 输入缓冲区 */
   tmpidxptr *item
   length int		/* 缓冲区当前有多少个数 */
   offset int	/* 缓冲区读到了文件的哪个位置 */
   idx int		/* 缓冲区的指针 */
 }

 /*
  * the tmp index struct for out sort
  */
type tmpIndex struct {
  k int
  bufs [K]tmpIndexBuf
  ls [K]int
}

func (idx *tmpIndexBuf) addIndexItem(d *dictionary, key string, docid int) {
  wrdid := d.getWordValue(key)
  if (offset + 8 >= maxbufsize {
    d.writer.Write(index.buf)
    offset = 0
  }

  wrdbuf = bytes.NewBuffer([]byte{})
  docbuf = bytes.NewBuffer([]byte{})
  binary.Write(wrdbuf, binary.BigEndian, wrdid)
  binary.Write(docbuf, binary.BigEndian, docid)
  copy(idx.buf[offset:], wrdbuf)
  copy(idx.buf[offset:], docbuf)
  offset += 8
}

 /*
  * the mian enter of algorithm of quicksort
  */
func (idx *tmpIndexBuf)quickSort(s, t int) {
  m := idx.split()
  for s > t {
    idx.quickSort(idx, s, m-1)
    idx.quickSort(idx, m+1, t)
  }
}
/*
 * the split part of algorithm of quicksort
 */
func (idx *tmpIndexBuf)split(idx []int, s, t int) int {
  for i, j:= s; i < t; i++ {
    if idx.Less(i, t) {
      idx.Swap(i, j)
      j++
    }
  }
  idx.Swap(j, t)
  return j
}
func (idx *tmpIndexBuf) Swap(i, j int) {
  idx.tmpidxptr[i], idx.tmpidxptr[j] = idx.tmpidxptr[j], idx.tmpidxptr[i]
}
func (idx *tmpIndexBuf) Less(i, j int) bool {
  if idx.tmpidxptr[i].wordid < idx.tmpidxptr[j].wordid {
    return true
  }

  if idx.tmpidxptr[i].docid < idx.tmpidxptr[j].docid {
    return true
  }

  return false
}

func (idx *tmpIndex) merge() {
  idx.buildLoseTree()
}

func (idx *tmpIndex) buildLoseTree() {
  for i := 0; i < k; i++ {
    idx.ls[i] = -1
  }
  for i := 0; i < k; i++ {
    idx.adjust(i)
  }
}

func (idx *tmpIndex) adjust(s int) {
  t := (idx.k + s) >> 1
  for t > 0 {
    if s == -1 {
      break
    }
    if (t == -1 || idx.bufs[s].tmpidxptr[offset] > idx.bufs[t].tmpidxptr[offset]) {
      s, ls[t] = ls[t], s
    }
    t >>= 1
  }
  ls[0] = s
}
