/*
** A complete page cache is an instance of this structure.  Every
** entry in the cache holds a single page of the database file.  The
** btree layer only operates on the cached copy of the database pages.
**
** A page cache entry is "clean" if it exactly matches what is currently
** on disk.  A page is "dirty" if it has been modified and needs to be
** persisted to disk.
**
** pDirty, pDirtyTail, pSynced:
**   All dirty pages are linked into the doubly linked list using
**   PgHdr.pDirtyNext and pDirtyPrev. The list is maintained in LRU order
**   such that p was added to the list more recently than p->pDirtyNext.
**   PCache.pDirty points to the first (newest) element in the list and
**   pDirtyTail to the last (oldest).
*/

struct PCache {
  PgHdr *pDirty, *pDirtyTail;         /* List of dirty pages in LRU order */
  PgHdr *pSynced;                     /* Last synced page in dirty page list */
  int nRefSum;                        /* Sum of ref counts over all pages */
  int szCache;                        /* Configured cache size */
  int szSpill;                        /* Size before spilling occurs */
  int szPage;                         /* Size of every page in this cache */
  int szExtra;                        /* Size of extra space for each page */

  /* Hash table of all pages. The following variables may only be accessed
  ** when the accessor is holding the PGroup mutex.
  */
  unsigned int nRecyclable;           /* Number of pages in the LRU list */
  unsigned int nPage;                 /* Total number of pages in apHash */
  unsigned int nHash;                 /* Number of slots in apHash[] */
  PgHdr1 **apHash;                    /* Hash table for fast lookup by key */
  PgHdr1 *pFree;                      /* List of unused pcache-local pages */
  void *pBulk;                        /* Bulk memory used by pcache-local */
};

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
  PgHdr *pDirtyNext;             /* Next element in list of dirty pages */
  PgHdr *pDirtyPrev;             /* Previous element in list of dirty pages */
};

/* Initialize and shutdown the page cache subsystem */
int cacheInitialize(void);
void cacheShutdown(void);

/* Page cache buffer management:
** These routines implement SQLITE_CONFIG_PAGECACHE.
*/
void cacheBufferSetup(void *, int sz, int n);

/* Create a new pager cache.
** Under memory stress, invoke xStress to try to make pages clean.
** Only clean and unpinned pages can be reclaimed.
*/
int cacheOpen(
  int szPage,                    /* Size of every page */
  int szExtra,                   /* Extra space associated with each page */
  int bPurgeable,                /* True if pages are on backing store */
  int (*xStress)(void*, PgHdr*), /* Call to try to make pages clean */
  void *pStress,                 /* Argument to xStress */
  PCache *pToInit                /* Preallocated space for the PCache */
);

/* Modify the page-size after the cache has been created. */
int cacheSetPageSize(PCache *, int);

/* Return the size in bytes of a PCache object.  Used to preallocate
** storage space.
*/
int cacheSize(void);

/* One release per successful fetch.  Page is pinned until released.
** Reference counted.
*/
sqlite3_pcache_page *cacheFetch(PCache*, Pgno, int createFlag);
int cacheFetchStress(PCache*, Pgno, sqlite3_pcache_page**);
PgHdr *cacheFetchFinish(PCache*, Pgno, sqlite3_pcache_page *pPage);
void cacheRelease(PgHdr*);

void cacheDrop(PgHdr*);         /* Remove page from cache */
void cacheMakeDirty(PgHdr*);    /* Make sure page is marked dirty */
void cacheMakeClean(PgHdr*);    /* Mark a single page as clean */
void cacheCleanAll(PCache*);    /* Mark all dirty list pages as clean */

/* Get a list of all dirty pages in the cache, sorted by page number */
PgHdr *cacheDirtyList(PCache*);

/* Reset and close the cache object */
void cacheClose(PCache*);

/* Clear flags from pages of the page cache */
void cacheClearSyncFlags(PCache *);

/* Discard the contents of the cache */
void cacheClear(PCache*);
/* Free up as much memory as possible from the page cache */
void cacheShrink(PCache*);

/* Return the header size */
int sqlite3HeaderSizePcache(void);

/* Number of dirty pages as a percentage of the configured cache size */
int cachePercentDirty(PCache*);
