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
**   such that p was added to the list more recently than p.pDirtyNext.
**   PCache.pDirty points to the first (newest) element in the list and
**   pDirtyTail to the last (oldest).
*/

struct PCache {
  szPage int                         /* Size of database content section */
  szExtra int                        /* sizeof(MemPage)+sizeof(PgHdr) */
  szAlloc int                     /* Total size of one pcache line */
  nMin uint32                  /* Minimum number of pages reserved */
  nMax uint32                  /* Configured "cache_size" value */
  pBulk *[]byte

  /* Hash table of all pages. The following variables may only be accessed
  ** when the accessor is holding the PGroup mutex.
  */
  nPage uint32                 /* Total number of pages in apHash */
  nHash uint32                /* Number of slots in apHash[] */
  apHash **PgHdr                    /* Hash table for fast lookup by key */
  pNext *PgHdr                     /* Next in hash table chain */
  iKey uint32                  /* Key value (page number) */
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
type PgHdr struct {
  pData *[]byte                   /* Page data */
  pExtra *[]byte                  /* Extra content */
  pCache *PCache              /* PRIVATE: Cache that owns this page */
  pDirty *PgHdr                 /* Transient list of dirty sorted by pgno */
  pPager *PgHdr                  /* The pager this page is part of */
  iKey uint32                     /* Page number for this page */
  pDirtyNext *PgHdr             /* Next element in list of dirty pages */
  pDirtyPrev *PgHdr             /* Previous element in list of dirty pages */
  pLruNext *PgHdr             /* Next in LRU list of unpinned pages */
  pLruPrev *PgHdr              /* Previous in LRU list of unpinned pages */
}

void cacheMakeDirty(PgHdr*);    /* Make sure page is marked dirty */
void cacheMakeClean(PgHdr*);    /* Mark a single page as clean */
void cacheCleanAll(PCache*);    /* Mark all dirty list pages as clean */

/*
** Implementation of the Create method.
**
** Allocate a new cache.
*/
func (pCache *PCache) Create(int szPage, int szExtra) {
  pCache.szPage = szPage;
  pCache.szExtra = szExtra;
  pCache.szAlloc = szPage + szExtra + ROUND8(sizeof(PgHdr));
  // pcache1EnterMutex(pGroup);
  pCache.ResizeHash();
  // pcache1LeaveMutex(pGroup);
  if( pCache.nHash==0 ){
    pCache.Destroy();
    pCache = 0;
  }
  pCache.InitBulk()
}

/*
** Try to initialize the pCache.pFree and pCache.pBulk fields.  Return
** true if pCache.pFree ends up containing one or more free pages.
*/
func (pCache *PCache) InitBulk() *[]byte {
  /* Do not bother with a bulk allocation if the cache size very small */
  szBulk := pCache.nInitPage>0 ? pCache.szAlloc * pCache.nInitPage : pCache.szAlloc * 1024

  zBulk := pCache.pBulk = make([]byte, szBulk)
  int nBulk = szBulk/pCache.szAlloc
  for --nBulk {
    PgHdr *pX = (PgHdr*)&zBulk[pCache.szPage];
    pX.pBuf = zBulk;
    pX.pExtra = &pX[1];
    pX.pNext = pCache.pFree;
    pCache.pFree = pX;
    zBulk += pCache.szAlloc;
  }
  return pCache.pFree!=0;
}


/*
** Implementation of the Destroy method.
**
** Destroy a cache allocated using Create().
*/
func (pCache *PCache) Destroy(){
  // if( pCache.nPage ) pcache1TruncateUnsafe(pCache, 0);
  // free(pCache.apHash);
  // free(pBulk)
  // free(pCache);
}

func (pCache *PCache) FetchPage(iKey uint32) *PgHdr {

  /* Step 1: Search the hash table for an existing entry. */
  pPage := pCache.apHash[iKey % pCache.nHash];
  for pPage && pPage.iKey!=iKey {
    pPage = pPage.pNext;
  }

  /* Step 2: If the page was found in the hash table, then return it.
  ** If the page was not in the hash table continue with
  ** subsequent steps to try to create the page. */
  if pPage != nil {
      return pPage
  }
  /* Steps 3 if page num is nearly full resize the hash*/
  if pCache.nPage>=pCache.nHash {
    pCache.ResizeHash()
  }
  /* Step 4. Try to recycle a page. */
  if pCache.nPage+1>=pCache.nMax /*|| pcache1UnderMemoryPressure(pCache)*/ {
    pPage = pGroup.lru.pLruPrev
    pCache.RemoveFromHash(pPage)
  }
  /* Step 5. If a usable page buffer has still not been found,
  ** attempt to allocate a new one.
  */
  if pPage == nil {
    pPage = pCache.AllocPage(pCache, createFlag==1);
  }

  if pPage != nil {
    h := iKey % pCache.nHash;
    pCache.nPage++;
    pPage.iKey = iKey;
    pPage.pNext = pCache.apHash[h];
    pPage.pCache = pCache;
    pPage.pLruPrev = 0;
    pPage.pLruNext = 0;
    pPage.isPinned = 1;
    *(void **)pPage.page.pExtra = 0;
    pCache.apHash[h] = pPage;
    if( iKey>pCache.iMaxKey ){
      pCache.iMaxKey = iKey;
    }
  }
  return pPage;
}

/*
** Allocate a new page object initially associated with cache pCache.
*/
func (pCache *PCache) AllocPage(){
  if pCache.pFree /*|| (pCache.nPage==0 && pcache1InitBulk(pCache))*/{
    p = pCache.pFree
    pCache.pFree = p.pNext
    p.pNext = 0
  } else {
    pPg := make([]byte, pCache.szAlloc)
    p = (PgHdr1 *)&((u8 *)pPg)[pCache.szPage]
    if( pPg==0 ) return 0
    p.page.pBuf = pPg
    p.page.pExtra = &p[1]
    p.isBulkLocal = 0
  }
  return p
}

/*
** Free a page object allocated by pcache1AllocPage().
*/
func (pCache *PCache) FreePage(PgHdr *p){

  // if( p.isBulkLocal ){
  p.pNext = pCache.pFree;
  pCache.pFree = p;
}

/*
** This function is used to resize the hash table used by the cache passed
** as the first argument.
**
** The PCache mutex must be held when this function is called.
*/
func (pCache *PCache) ResizeHash(){

  nNew := pCache.nHash*2;
  if( nNew<256 ){
    nNew = 256;
  }

  apNew := make([]*PgHdr, nNew);

  for i:=0; i<pCache.nHash; i++{
    pCurPg := pCache.apHash[i];
    for pCurPg != nil {
      h := pCurPg.iKey % nNew;
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

  h := pPage.iKey % pCache.nHash
  p := &pCache.apHash[h]
  for p; (*p)!=pPage; p=&(*p).pNext);
  *p = (*p).pNext;

  pCache.nPage--;
  pCache.FreePage(pPage);
}

/* Allowed values for second argument to ManageDirtyList() */
const PCACHE_DIRTYLIST_REMOVE   1    /* Remove pPage from dirty list */
const PCACHE_DIRTYLIST_ADD      2    /* Add pPage to the dirty list */
const PCACHE_DIRTYLIST_FRONT    3    /* Move pPage to the front of the list */

/*
** Manage pPage's participation on the dirty list.  Bits of the addRemove
** argument determines what operation to do.  The 0x01 bit means first
** remove pPage from the dirty list.  The 0x02 means add pPage back to
** the dirty list.  Doing both moves pPage to the front of the dirty list.
*/
func (pCache *PCache) ManageDirtyList(pPage *PgHdr, addRemove uint8){

  if addRemove & PCACHE_DIRTYLIST_REMOVE {

    /* Update the PCache.pSynced variable if necessary. */
    // if( p.pSynced==pPage ){
    //   p.pSynced = pPage.pDirtyPrev;
    // }

    if pPage.pDirtyNext != nil {
      pPage.pDirtyNext.pDirtyPrev = pPage.pDirtyPrev;
    }else{
      pCache.pDirtyTail = pPage.pDirtyPrev;
    }
    if pPage.pDirtyPrev != nil {
      pPage.pDirtyPrev.pDirtyNext = pPage.pDirtyNext;
    }else{
      /* If there are now no dirty pages in the cache, set eCreate to 2.
      ** This is an optimization that allows sqlite3PcacheFetch() to skip
      ** searching for a dirty page to eject from the cache when it might
      ** otherwise have to.  */
      pCache.pDirty = pPage.pDirtyNext;
    }
    pPage.pDirtyNext = 0;
    pPage.pDirtyPrev = 0;
  }
  if( addRemove & PCACHE_DIRTYLIST_ADD ){
    pPage.pDirtyNext = p.pDirty;
    if( pPage.pDirtyNext ){
      pPage.pDirtyNext.pDirtyPrev = pPage;
    }else{
      pCache.pDirtyTail = pPage;
    }
    pCache.pDirty = pPage;
  }
}
