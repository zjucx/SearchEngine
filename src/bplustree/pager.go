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
  "syscall"
  "unsafe"
  "os"
  "fmt"
)

type PgHead struct {
  flag  uint8
  ncell int
  nfree int                   /* free bytes in current page */
  pgno  int                  /* page number */
  ppgno int                  /* parent page number */
  maxkey int                 /* max key int current page */
}

type Pager struct{
  f *os.File              /* Number of mmap pages currently outstanding */
  szPage int               /* Number of bytes in a page */
  numPage int                /* Maximum allowed size of the database */
  dbName string           /* Name of the database file */
  pCache *PCache;            /* Pointer to page cache object */
};

/* Open and close a Pager connection. */
func (p *Pager) Open(dbName string, szPage int) {
  p.szPage = szPage
  p.dbName = dbName

  f, err := os.OpenFile(dbName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
  if err != nil {
		fmt.Println(err)
	}
  p.f = f

  fi, err := os.Stat(dbName)
  numPage := fi.Size()/int64(p.szPage)

  pCache := &PCache{}
  pCache.Create(szPage)
  p.pCache = pCache

  if numPage == 0 {
    pg0 := p.Fetch(0)
    pg1 := p.Fetch(1)
    pg0.WritePageHeader(0, 0, 0, 0, 0, 0)
    p.Write(pg0)
    p.Sync()
    numPage = 2
  }
  p.numPage = int(numPage)
}

func (p *Pager) Destroy() {
  if p.f != nil {
    p.f.Close()
  }
  p.pCache.Destroy()
}

func (p *Pager) Fetch(pgno int) (*PgHdr){
  if pgno > p.numPage {
    // log
    return nil
  }
  return p.pCache.FetchPage(pgno)
}

func (p *Pager) Read(pgno int) (pPg *PgHdr){
  pPg = p.Fetch(pgno)
  if pPg == nil {
    return nil
  }
  szPage := p.pCache.szPage
  n, err := p.f.ReadAt(pPg.pBuf[:szPage], pgno * szPage)
  if err != nil {
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
  pgHead := (*PgHead)(pPg.GetPageHeader())
  //n, err := p.f.WriteAt(pPg.pBuf[:szPage], (pPg.pgno-1) * szPage)

  /* Update the database size and return. */
  if(p.numPage < pgHead.pgno){
    p.numPage = pgHead.pgno;
  }
  return 0
}

/*
** Sync the database file to disk. This is a no-op for in-memory databases
** or pages with the Pager.noSync flag set.
*/
func (p *Pager) Sync(){
  // sync file func Fdatasync(fd int) (err error)
  err := syscall.Fdatasync(p.f.Fd())
  if err != nil {
    // log
    return
  }
  // make cache clear
  p.pCache.MakeCleanAll();
}

func (pgHdr *PgHdr) GetPageHeader() unsafe.Pointer {
  return unsafe.Pointer(pgHdr.pBuf)
}

func (pgHdr *PgHdr) WritePageHeader(flag uint8, ncell, nfree int,
  pgno, ppgno, maxkey int) {
  pgHead := (*PgHead)(unsafe.Pointer(pgHdr.pBuf))
  pgHead.flag = flag
  pgHead.ncell = ncell
  pgHead.nfree = nfree
  pgHead.pgno = pgno
  pgHead.ppgno = ppgno
  pgHead.maxkey = maxkey
}
