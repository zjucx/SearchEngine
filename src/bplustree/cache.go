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
  PgHdr *pFree;                       /* List of unused pcache-local pages */
  int szCache;                        /* Configured cache size */
  int szSpill;                        /* Size before spilling occurs */
  int szPage;                         /* Size of every page in this cache */
  int szExtra;                        /* Size of extra space for each page */
  pBulk *unsafe.Poiter

  /* Hash table of all pages. The following variables may only be accessed
  ** when the accessor is holding the PGroup mutex.
  */
  unsigned int nRecyclable;           /* Number of pages in the LRU list */
  unsigned int nPage;                 /* Total number of pages in apHash */
  unsigned int nHash;                 /* Number of slots in apHash[] */
  apHash **PgHdr                     /* Hash table for fast lookup by key */
  PgHdr *pNext;                       /* Next in hash table chain */
  unsigned int iKey;                  /* Key value (page number) */
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



/*
** Implementation of the Create method.
**
** Allocate a new cache.
*/
func (pCache *PCache) Create(int szPage, int szExtra) {
  pCache->szPage = szPage;
  pCache->szExtra = szExtra;
  pCache->szAlloc = szPage + szExtra + ROUND8(sizeof(PgHdr));
  // pcache1EnterMutex(pGroup);
  pCache.ResizeHash();
  // pcache1LeaveMutex(pGroup);
  if( pCache->nHash==0 ){
    pCache.Destroy();
    pCache = 0;
  }
  pCache.InitBulk()
}

/*
** Try to initialize the pCache->pFree and pCache->pBulk fields.  Return
** true if pCache->pFree ends up containing one or more free pages.
*/
func (pCache *PCache) InitBulk() unsafe.Pointer {
  if( pCache.nInitPage==0 ) return
  /* Do not bother with a bulk allocation if the cache size very small */
  if( pCache.nInitPage>0 ){
    szBulk := pCache->szAlloc * pCache.nInitPage;
  }else{
    szBulk = -1024 * (i64)pcache1.nInitPage;
  }

  zBulk = pCache->pBulk = new( szBulk );
  if( zBulk ){
    int nBulk = sqlite3MallocSize(zBulk)/pCache->szAlloc;
    do{
      PgHdr *pX = (PgHdr*)&zBulk[pCache->szPage];
      pX->page.pBuf = zBulk;
      pX->page.pExtra = &pX[1];
      pX->pNext = pCache->pFree;
      pCache->pFree = pX;
      zBulk += pCache->szAlloc;
    }while( --nBulk );
  }
  return pCache->pFree!=0;
}


/*
** Implementation of the Destroy method.
**
** Destroy a cache allocated using Create().
*/
func (pCache *PCache) Destroy(){
  if( pCache->nPage ) pcache1TruncateUnsafe(pCache, 0);
  free(pCache->apHash);
  free(pBulk)
  free(pCache);
}

/*
** This function is used to resize the hash table used by the cache passed
** as the first argument.
**
** The PCache mutex must be held when this function is called.
*/
func (pCache *PCache) ResizeHash(){

  nNew := pCache->nHash*2;
  if( nNew<256 ){
    nNew = 256;
  }

  apNew := (PgHdr **)new(sizeof(PgHdr *)*nNew);

  for i:=0; i<pCache->nHash; i++{
    pCurPg := pCache->apHash[i];
    for pCurPg != nil {
      h := pCurPg->iKey % nNew;
      pNewPg := apNew[h]

      apNew[h] = pCurPg
      pCurPg = pCurPg.pNext
      apNew[h].pNext = pNewPg
    }
  }
  free(pCacheapHash);
  pCache.apHash = apNew;
  pCache.nHash = nNew;
}

/*
** Remove the page supplied as an argument from the hash table
** (PCache1.apHash structure) that it is currently stored in.
** Also free the page if freePage is true.
**
*/
func (pCache *PCache) RemoveFromHash(pPage *PgHdr) {

  h := pPage->iKey % pCache->nHash
  p := &pCache->apHash[h]
  for p; (*p)!=pPage; p=&(*p)->pNext);
  *p = (*p)->pNext;

  pCache->nPage--;
  pCache.FreePage(pPage);
}
