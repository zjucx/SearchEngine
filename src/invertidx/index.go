package invertidx

import (
  "fmt"
  "unsafe"
  "bufio"
  "os"
  //"encoding/binary"
)

const (
	tindex     uint8 = 0
	tdata      uint8 = 1
)

const K = 64
const maxbufsize = 4096

type Item struct {
  docid     uint32
  wordid    uint32
}

/*
 * the tmp index file(.tmp) 's data orgnized by
 */
 type IndexBuf struct {
   buf []byte		/* 输入缓冲区 */
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

func NewIdxBuf(filename string) IndexBuf{
  idxBuf := IndexBuf{
    filename : filename,
    buf : make([]byte, maxbufsize),
    length : 0,
    offset : 0,
    idx : 0,
  }
  return idxBuf
}

func (idx *IndexBuf) AddIndexItem(docid, wordid int) {
  item := Item {
    docid : uint32(docid),
    wordid : uint32(wordid),
  }

  l := unsafe.Sizeof(item)
  pb := (*[8]byte)(unsafe.Pointer(&item))

  if idx.offset + int(l) >= maxbufsize {
    idx.Flush()
  }
  //binary.Write(&idx.buf[idx.offset], binary.BigEndian, item)
  bufptr := (*[maxbufsize]byte)(unsafe.Pointer(&idx.buf[idx.offset]))
  copy((*bufptr)[:l], (*pb)[:l])
  idx.offset += int(l)
  //tmp := (*[maxbufsize]byte)(unsafe.Pointer(&idx.buf[0]))
  //fmt.Println((*tmp)[:idx.offset])
}

func (idx *IndexBuf) Flush(){
  //open file and write to file
  f, err := OpenFile(idx.filename)
  CheckErr(err)
  bw := bufio.NewWriter(f)
  idx.quickSort(0, idx.offset/int(unsafe.Sizeof(&Item{}))-1)
  bw.Write(idx.buf[:idx.offset])
  bw.Flush()
  // f.Write(idx.buf)
  // f.Close()
  //reset var
  idx.offset = 0
  idx.length = 0
  // idx.buf = make([]byte, maxbufsize)
}
 /*
  * the mian enter of algorithm of quicksort
  */
func (idx *IndexBuf)quickSort(s, t int) {
  if t < 0 {return}
  m := idx.split(s, t)
  for s > t {
    idx.quickSort(s, m-1)
    idx.quickSort(m+1, t)
  }
}
/*
 * the split part of algorithm of quicksort
 */
func (idx *IndexBuf)split(s, t int) int {
  var i, j int
  for i, j = s, s; i < t; i++ {
    if idx.less(i, t) {
      idx.swap(i, j)
      j++
    }
  }
  idx.swap(j, t)
  return j
}
func (idx *IndexBuf) swap(i, j int) {
  ptr := (*[maxbufsize/unsafe.Sizeof(&Item{})]Item)(unsafe.Pointer(&idx.buf[0]))
  (*ptr)[i].wordid, (*ptr)[j].wordid = (*ptr)[j].wordid, (*ptr)[i].wordid
  (*ptr)[i].docid, (*ptr)[j].docid = (*ptr)[j].docid, (*ptr)[i].docid
}
func (idx *IndexBuf) less(i, j int) bool {
  ptr := (*[maxbufsize/unsafe.Sizeof(&Item{})]Item)(unsafe.Pointer(&idx.buf[0]))
  if (*ptr)[i].docid < (*ptr)[j].docid {
    return true
  }
  if (*ptr)[i].wordid < (*ptr)[j].wordid {
    return true
  }
  return false
}

func (idx *Index)readDataFromFile(filename string, offset, bufIdx int) int{
  fi, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  bfRd := bufio.NewReader(fi)
  fi.Seek(int64(offset), 0)
  bytes, err := bfRd.Read(idx.bufs[bufIdx].buf)
  idx.bufs[bufIdx].length = bytes / int(unsafe.Sizeof(&Item{}))
  return bytes
}

func (idx *Index) sortIndexFile(filename string) {
  f, e := os.Stat(filename)
	if e != nil {
		return
	}
  //file size
  filesize := f.Size()
  //numFile = 1;
  //the number of bufs to load the file
  runNum := int(filesize / maxbufsize)
  //leftBuf := filesize % maxbufsize*K
  fi, err := os.Open(filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  bfRd := bufio.NewReader(fi)

  //file name format  loopidx+mergeidx.tindex
  //the frist loop will generite file which the name prefix is 1i.tindex
  //the second loop is 2i.tindex
  prefix := 1;

  for runNum > 0 {
    //N bufs can merge to one file
    runNum = runNum / K
    if runNum % K > 0{
      runNum++
    }

    for i := 0; i < runNum; i++ {
      if i == runNum-1 && runNum % K > 0 {
        idx.k = runNum % K
      } else {
        idx.k = K
      }

      // read buf from file if numFile == 1 read from original file
      // else from tmp index files
      for j := 0; j < idx.k; j++ {
        var bytes int
        if i == 1 {
          bytes, err := bfRd.Read(idx.bufs[j].buf)
          CheckErr(err)
          idx.bufs[j].length = bytes / int(unsafe.Sizeof(&Item{}))
        } else {
          filename := fmt.Sprintf("%d%d.tindex", prefix-1, i*K+j)
          bytes = idx.readDataFromFile(filename, 0, j)
        }
        idx.bufs[j].offset = bytes
        idx.bufs[j].idx = 0
      }
      idx.merge(prefix, i)
    }
    prefix++
  }
}

func (idx *Index) merge(prefix, curNum int) {
  idx.buildLoseTree(idx.k)
  //var filename string
  filename := fmt.Sprintf("%d%d.tindex", prefix, curNum)
  //snprintf(filename, 100, "%s%d.dat", output_prefix, n_merge)
  fo, _ := os.Create(filename)  //创建文件
  bw := bufio.NewWriter(fo)

  fmt.Println("file is not exist!");
  k := idx.k
  for k > 0 {
    mr := idx.bufs[idx.ls[0]]
    idx.bufo.buf[idx.bufo.idx] = mr.buf[mr.idx]
    idx.bufo.idx++
    mr.idx++
    //output buf is full
    if idx.bufo.idx == maxbufsize {
      idx.bufo.idx = 0
      bw.Write(idx.bufo.buf)
      //write to file
    }

    //input buf is full
    if mr.idx == mr.length {
      //read data from file until file EOF
      filename := fmt.Sprintf("%d%d.tindex", prefix-1, curNum*K+idx.ls[0])
      bytes := idx.readDataFromFile(filename, mr.offset, idx.ls[0])
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

func (idx *Index) buildLoseTree(k int) {
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
    if (t == -1 || (*(*Item)(unsafe.Pointer(&idx.bufs[s].buf[idx.bufs[s].offset]))).wordid > (*(*Item)(unsafe.Pointer(&idx.bufs[t].buf[idx.bufs[t].offset]))).wordid) {
      s, idx.ls[t] = idx.ls[t], s
    }
    t >>= 1
  }
  idx.ls[0] = s
}
