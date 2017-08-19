/*
** Each btree pages is divided into three sections:  The header, the
** cell pointer array, and the cell content area.  Page 1 also has a 100-byte
** file header that occurs before the page header.
**
**      |----------------|
**      | file header    |   100 bytes.  Page 1 only.
**      |----------------|
**      | page header    |   8 bytes for leaves.  12 bytes for interior nodes
**      |----------------|
**      | cell pointer   |   |  2 bytes per cell.  Sorted order.
**      | array          |   |  Grows downward
**      |                |   v
**      |----------------|
**      | unallocated    |
**      | space          |
**      |----------------|   ^  Grows upwards
**      | cell content   |   |  Arbitrary order interspersed with freeblocks.
**      | area           |   |  and free space fragments.
**      |----------------|
**
** The page headers looks like this:
**
**   OFFSET   SIZE     DESCRIPTION
**      0       1      Flags. 1: interpage, 2: leafpage, 4: overflowpage
**      1       2      byte offset to the first freeblock
**      3       2      number of cells on this page
**      5       2      first byte of the cell content area
**      7       1      number of fragmented free bytes
**      8       4      Right child (the Ptr(N) value).  Omitted on leaves.
*/
package bplustree

import(
  "unsafe"
  "os"
  "fmt"
  "encoding/binary"
  "bytes"
)

type PgHead struct {
  flag  uint8
  ncell uint16
  nfree uint16                   /* free bytes in current page */
  pgno  uint32                  /* page number */
  ppgno uint32                  /* parent page number */
  maxkey uint32                 /* max key int current page */
}

type Pager struct{
  file *os.File              /* Number of mmap pages currently outstanding */
  szPage  uint16               /* Number of bytes in a page */
  numPage uint32                /* Maximum allowed size of the database */
  dbName  string           /* Name of the database file */
  pCache  *PCache;            /* Pointer to page cache object */
};

/* Open and close a Pager connection. */
func (p *Pager) Create(dbName string, szPage uint16) {
  p.szPage = szPage
  p.dbName = dbName

  // open db file a db file is associated with a pager
  file, err := os.OpenFile(dbName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if err != nil {
		fmt.Println(err)
	}
  p.file = file

  fi, err := os.Stat(dbName)
  numPage := fi.Size()/int64(p.szPage)

  // init cache model
  p.pCache = &PCache{}
  p.pCache.Create(szPage)

  // get page from cache
  var pg0 *PgHdr

  if numPage == 0 {
    pg0 = p.Fetch(0)
    pg0.WriteHeader(1, 1, 1, 1, 1, 1)
    pg0.GetHeader()
    p.Write(pg0)
    p.Sync()
    numPage = 1
  } else {
    pg0 = p.Read(0)
  }
  p.numPage = uint32(numPage)
}

func (p *Pager) Close() {
  if p.file != nil {
    p.file.Close()
  }
  p.pCache.Destroy()
}

func (p *Pager) Read(pgno uint32) (pPg *PgHdr){
  pPg = p.Fetch(pgno)
  if pPg == nil {
    fmt.Println("get page from cache error!")
    return nil
  }

  szPage := p.pCache.szPage
  pBulk := *(*[]byte)(unsafe.Pointer(pPg.pBulk))
  n, err := p.file.ReadAt(pBulk[:szPage], int64(pgno * uint32(szPage)))
  if err != nil || n != int(szPage) {
    return nil
  }
  return pPg
}

/* Operations on page references. */
func (p *Pager) Write(pPg *PgHdr) int {
  /* Mark the page that is about to be modified as dirty. */
  p.pCache.MakeDirty(pPg);
  //func Pwrite(fd int, p []byte, offset int64) (n int, err error)
  szPage := p.pCache.szPage
  pgHead := (*PgHead)(pPg.GetHeader())
  pBulk := *(*[]byte)(unsafe.Pointer(pPg.pBulk))
  println(len(pBulk))

  n, err := p.file.WriteAt(pBulk[:szPage], int64(pgHead.pgno * uint32(szPage) - uint32(szPage)))
  if err != nil || n != int(szPage) {
    return 0
  }
  //n, err := p.file.WriteAt(pPg.pBulk[:szPage], (pPg.pgno-1) * szPage)

  /* Update the database size and return. */
  if(p.numPage < pgHead.pgno){
    p.numPage = pgHead.pgno;
  }
  return n
}

/*
** Sync the database file to disk. This is a no-op for in-memory databases
** or pages with the Pager.noSync flag set.
*/
func (p *Pager) Sync(){
  // sync file func Fdatasync(fd int) (err error)

  // make cache clear
  p.pCache.MakeCleanAll();
}

func (p *Pager) Fetch(pgno uint32) (*PgHdr){
  return p.pCache.FetchPage(pgno)
}

/*
type Head struct {
    Cmd byte
    Version byte
    Magic   uint16
    Reserve byte
    HeadLen byte
    BodyLen uint16
}

func NewHead(buf []byte)*Head{
    head := new(Head)
    head.Cmd     = buf[0]
    head.Version = buf[1]
    head.Magic   = binary.BigEndian.Uint16(buf[2:4])
    head.Reserve = buf[4]
    head.HeadLen = buf[5]
    head.BodyLen = binary.BigEndian.Uint16(buf[6:8])
    return head
}
*/

func (pgHdr *PgHdr) GetHeader() *PgHead {
  buf := *(*[]byte)(unsafe.Pointer(pgHdr.pBulk))
  fmt.Printf("len=%d cap=%d slice=%v\n",len(buf),cap(buf),buf)
  pgHead := &PgHead{
    flag: buf[0],
    ncell: binary.LittleEndian.Uint16(buf[1:3]),
    nfree: binary.LittleEndian.Uint16(buf[3:5]),
    pgno: binary.LittleEndian.Uint32(buf[5:9]),
    ppgno: binary.LittleEndian.Uint32(buf[9:13]),
    maxkey: binary.LittleEndian.Uint32(buf[13:17]),
  }
  fmt.Printf("len=%d cap=%d\n",pgHead.flag, pgHead.nfree)
  return pgHead
}

func (pgHdr *PgHdr) WriteHeader(flag uint8, ncell, nfree uint16,
  pgno, ppgno, maxkey uint32) {
  buf := *(*[]byte)(unsafe.Pointer(pgHdr.pBulk))

  copy(buf, ToBytes(flag))
  copy(buf[1:], ToBytes(ncell))
  copy(buf[3:], ToBytes(nfree))
  copy(buf[5:], ToBytes(pgno))
  copy(buf[9:], ToBytes(ppgno))
  copy(buf[13:], ToBytes(maxkey))
  fmt.Printf("len=%d cap=%d slice=%v\n",len(buf),cap(buf),buf)
}

func ToBytes(data interface{}) []byte {
  buf := new(bytes.Buffer)
  switch data.(type){
  case uint8:
    err := binary.Write(buf, binary.LittleEndian, uint8(data.(uint8)))
    if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
  case uint16:
    err := binary.Write(buf, binary.LittleEndian, int16(data.(uint16)))
    if err != nil {
			fmt.Println("binary.Write failed:", err)
		}
  case uint32:
    err := binary.Write(buf, binary.LittleEndian, int32(data.(uint32)))
    if err != nil {
      fmt.Println("binary.Write failed:", err)
    }
  }
  return buf.Bytes()
}
