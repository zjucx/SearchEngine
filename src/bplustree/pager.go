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

type PgHeader struct {
  flag uint8
  free uint16
  pgno uint32
  ppgno uint32
}

type Pager struct{
  f *File              /* Number of mmap pages currently outstanding */
  pageSize uint32               /* Number of bytes in a page */
  mxPgno uint32                /* Maximum allowed size of the database */
  fileName string           /* Name of the database file */
  PCache *pPCache;            /* Pointer to page cache object */
};

/* Open and close a Pager connection. */
func (p *Pager) Open(fileName string) {
  p.pageSize = 4096
  p.fileName = fileName

  f, err := OpenFile(filename, O_RDWR|O_APPEND|O_CREATE, 0660)
  if err != nil {
		fmt.Println(err)
	}
  p.f = f

  fi, err := os.Stat(path)
	if os.IsNotExist(err) {
		p.mxPgno = 0;
	}
  p.mxPgno = f.Size()/p.pageSize
  p.pCache.Open()
}

func (p *Pager) Close(filename string) {
  if p.f != nil {
    p.f.Close()
  }
  p.pCache.Close()
}
func (p *Pager)ReadPageHeader(pgno uint32) *PgHeader {

}
func (p *Pager)Shrink() {
  p.pCache.Shrink()
}

/* Operations on page references. */
int pagerWrite(DbPage*);
void pagerDontWrite(DbPage*);
int pagerMovepage(Pager*,DbPage*,Pgno,int);
int pagerPageRefcount(DbPage*);
void *pagerGetData(DbPage *);
void *pagerGetExtra(DbPage *);

/* Functions used to manage pager transactions and savepoints. */
void pagerPagecount(Pager*, int*);
int pagerCommit(Pager*,const char *zMaster, int);
int pagerSync(Pager *pPager, const char *zMaster);
