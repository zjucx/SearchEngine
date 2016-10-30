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
 type IndexBuf struct {
   buf [maxbufsize]byte		/* 输入缓冲区 */
   length int		/* 缓冲区当前有多少个数 */
   offset int	/* 缓冲区读到了文件的哪个位置 */
   idx int		/* 缓冲区的指针 */
   filename string
 }

 /*
  * the tmp index struct for out sort
  */
type Index struct {
  k int
  ls [K]int
  bufs [K]IndexBuf
  bufo IndexBuf
}

func (idx *IndexBuf) addIndexItem(docid, wordid int) {

  if (offset + size(item) >= maxbufsize {
    idx.Flush()
  }

  wrdbuf = bytes.NewBuffer([]byte{})
  docbuf = bytes.NewBuffer([]byte{})
  binary.Write(wrdbuf, binary.BigEndian, wrdid)
  binary.Write(docbuf, binary.BigEndian, docid)
  copy(idx.buf[offset:], wrdbuf)
  copy(idx.buf[offset:], docbuf)
  offset += size(item)
}

func (idx *IndexBuf) flush(){
  //open file and write to file
  f, err := OpenFile(idx.filename)
  CheckErr(err)
  bw := bufio.NewWriter(f)
  idx.quicksort(0, idx.offset/size(item))
  bw.Write(index.buf)
  bw.Flush()
  //reset var
  offset = 0
}
 /*
  * the mian enter of algorithm of quicksort
  */
func (idx *IndexBuf)quickSort(s, t int) {
  m := idx.split()
  for s > t {
    idx.quickSort(idx, s, m-1)
    idx.quickSort(idx, m+1, t)
  }
}
/*
 * the split part of algorithm of quicksort
 */
func (idx *IndexBuf)split(idx []int, s, t int) int {
  for i, j:= s; i < t; i++ {
    if idx.less(i, t) {
      idx.swap(i, j)
      j++
    }
  }
  idx.swap(j, t)
  return j
}
func (idx *IndexBuf) swap(i, j int) {
  ptr := (*item)(unsafe.Pointer(idx.buf))
  ptr[i], ptr[j] = ptr[j], ptr[i]
}
func (idx *IndexBuf) less(i, j int) bool {
  ptr := (*item)(unsafe.Pointer(idx.buf))
  if ptr[i].wordid < ptr[j].wordid {
    return true
  }

  if ptr[i].docid < ptr[j].docid {
    return true
  }

  return false
}

func (idx *Index)readDataFromFile(offset int) int{
  snprintf(filename, 20, "%s%d.dat", input_prefix, i*K+j);
  fi, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  bfR := bufio.NewReader(fi)
  _, _ := fi.Seek(offset, 0)
  bytes, err := bfRd.Read(idx.bufs[j].buf)
  idx.bufs[j].length = bytes / Size(int)
  return bytes
}

func (idx *Index) sortIndexFile(filename string) {
  f, e := os.Stat(filename)
	if e != nil {
		return
	}
  //file size
  filesize := f.Size()
  numFile = 1;
  //the number of bufs to load the file
  runNum := filesize / maxbufsize
  //leftBuf := filesize % maxbufsize*K
  fi, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  bfRd := bufio.NewReader(fi)

  for runNum {
    //N bufs can merge to one file
    runNum = runNum / N
    if runNum % N {
      runNum++
    }

    for i := 0; i < runNum; i++ {
      if i == runNum-1 && runNum % N {
        idx.k = runNum % N
      } else {
        idx.k = K
      }

      // read buf from file if numFile == 1 read from original file
      // else from tmp index files
      for j := 0; j < needMerge; j++ {
        if numFile == 1 {
          bytes, err := bfRd.Read(idx.bufs[j].buf)
          idx.bufs[j].length = bytes / Size(int)
        } else {
          bytes := idx.readDataFromFile(0)
        }
        idx.bufs[j].offset = bytes
        idx.bufs[j].idx = 0
      }
      merge(i)
    }
    numFile = 0
  }
}

func (idx *Index) merge(int curNum) {
  idx.buildLoseTree()
  var filename string
  snprintf(filename, 100, "%s%d.dat", output_prefix, n_merge)
  fo, _ := os.Create(filename)  //创建文件

  fmt.Println("file is not exist!");
  k := idx.k
  for k {
    mr := idx.bufs[idx.ls[0]]
    idx.bufo.buf[idx.bufo.idx++] = mr.buf[mr.idx++]

    //output buf is full
    if idx.bufo.idx == maxbufsize {
      idx.bufo.idx = 0
      //write to file
    }

    //input buf is full
    if mr.idx == mr.length {
      //read data from file until file EOF
      bytes := idx.readDataFromFile(mr.offset)
      if bytes == 0 {
        k--
      } else {
        mr.offset += bytes
        mr.idx = 0
      }
    }
    idx.adjust(idx.ls[0])
  }
  //write left data to file
  //bytes = write(output_fd, buffer, bp*4)
}

func (idx *Index) buildLoseTree() {
  for i := 0; i < k; i++ {
    idx.ls[i] = -1
  }
  for i := 0; i < k; i++ {
    idx.adjust(i)
  }
}

func (idx *Index) adjust(s int) {
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
