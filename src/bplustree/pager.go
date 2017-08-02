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

import(
  "syscall"
  "unsafe"
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
  f *File              /* Number of mmap pages currently outstanding */
  szPage uint32               /* Number of bytes in a page */
  numPage uint32                /* Maximum allowed size of the database */
  dbName string           /* Name of the database file */
  pCache *PCache;            /* Pointer to page cache object */
};

/* Open and close a Pager connection. */
func (p *Pager) Open(dbName string, szPage int) {
  p.szPage = szPage
  p.dbName = dbName

  f, err := OpenFile(dbName, O_RDWR|O_APPEND|O_CREATE, 0660)
  if err != nil {
		fmt.Println(err)
	}
  p.f = f

  fi, err := os.Stat(path)
  numPage = fi.Size()/p.szPage

  pCache := &PCache{}
  pCache.Create(szPage)
  p.pCache = pCache

  if numPage == 0 {
    pg0 := p.Fetch(0)
    pg1 := p.Fetch(1)
    pg0.WritePageHeader(0, 0, 0 , 0, 0, 0, 0)
    pg1.WritePageHeader(1, 0, 0 , 0, 0, 0, 0)
    p.Write(pg0)
    p.Write(pg1)
    p.Sync()
    numPage = 2
  }
  p.numPage = numPage
}

func (p *Pager) Close() {
  if p.f != nil {
    p.f.Close()
  }
  p.pCache.Close()
}

func (p *Pager) Shrink() {
  p.pCache.Shrink()
}

func (p *Pager) Fetch(pgno uint32) (*PgHdr){
  if pgno > p.numPage {
    // log
    return nil
  }
  return p.pCache.FetchPage(pgno)
}

func (p *Pager) Read(pgno uint32) (pPg *PgHdr){
  pPg := p.Fetch(pgno)
  if pPg == nil {
    return nil
  }
  szPage := p.pCache.szPage
  n, err = p.f.ReadAt(pPg.pBuf[:szPage], pPg.pgno * szPage)
  if err != nil {
    return nil
  }
  return pPg
}

/* Operations on page references. */
func (p *Pager) Write(pPg *PgHdr) (n int, err Error){
  /* Mark the page that is about to be modified as dirty. */
  p.pCache.MakeDirty(pPg);
  //func Pwrite(fd int, p []byte, offset int64) (n int, err error)
  szPage := p.pCache.szPage
  n, err = p.f.WriteAt(pPg.pBuf[:szPage], (pPg.pgno-1) * szPage)

  if err != nil || n != szPage {
    //log
    return n, err
  }
  /* Update the database size and return. */
  if( p.dbSize < pPg.pgno ){
    pPager.dbSize = pPg.pgno;
  }
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
  p.pCache.CleanAll();
}

func (pgHdr *PgHdr) GetCellPtr() unsafe.Pointer {
  return unsafe.Pointer(&(*PgHead)unsafe.Pointer(pgHdr.pBuf)[1])
}

func (pgHdr *PgHdr) GetPageHeader() *PgHead {
  return (*PgHead)unsafe.Pointer(pgHdr.pBuf)
}

func (pgHdr *PgHdr) WritePageHeader(flag uint8, ncell, nfree uint16,
  pgno, ppgno, maxkey uint32) {
  pgHead := (*PgHead)unsafe.Pointer(&pgHdr.pBuf[0])
  pgHead.flag = flag
  pgHead.ncell = ncell
  pgHead.nfree = nfree
  pgHead.pgno = pgno
  pgHead.ppgno = ppgno
  pgHead.maxkey = maxkey
}
