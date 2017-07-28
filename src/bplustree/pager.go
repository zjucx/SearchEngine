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

/*
** Every page in the cache is controlled by an instance of the following
** structure.
**
** A Page cache line looks like this:
**
**  --------------------------------------------------
**  |  database page content   |  PgHdr  |  MemPage  |
**  --------------------------------------------------
*/
struct PgHdr {
  void *pData;                   /* Page data */
  void *pExtra;                  /* Extra content */
  PCache *pCache;                /* PRIVATE: Cache that owns this page */
  PgHdr *pDirty;                 /* Transient list of dirty sorted by pgno */
  Pager *pPager;                 /* The pager this page is part of */
  Pgno pgno;                     /* Page number for this page */
  i16 nRef;                      /* Number of users of this page */
  PgHdr *pDirtyNext;             /* Next element in list of dirty pages */
  PgHdr *pDirtyPrev;             /* Previous element in list of dirty pages */
};

struct Pager {
  f *File              /* Number of mmap pages currently outstanding */
  PgHdr *pMmapFreelist;       /* List of free mmap page headers (pDirty) */
  /*
  ** End of the routinely-changing class members
  ***************************************************************************/

  u16 nExtra;                 /* Add this many bytes to each in-memory page */
  int pageSize;               /* Number of bytes in a page */
  Pgno mxPgno;                /* Maximum allowed size of the database */
  filename string           /* Name of the database file */
  PCache *pPCache;            /* Pointer to page cache object */
};

/* Open and close a Pager connection. */
func (p *Pager) Open(filename string) {
  p.filename = filename
  f, err := OpenFile(filename, O_RDWR|O_APPEND|O_CREATE, 0660)
  if err != nil {
		fmt.Println(err)
	}
  p.f = f
}

func (p *Pager) Close(filename string) {
  if p.f != nil {
    p.f.Close()
  }
}
int pagerReadFileheader(Pager*, int, unsigned char*);
void pagerShrink(Pager*);

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
