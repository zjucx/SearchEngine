/*
** A complete pgHdr cache is an instance of this structure.  Every
** entry in the cache holds a single pgHdr of the database file.  The
** btree layer only operates on the cached copy of the database pages.
**
** A pgHdr cache entry is "clean" if it exactly matches what is currently
** on disk.  A pgHdr is "dirty" if it has been modified and needs to be
** persisted to disk.
**
** pDirty, pDirtyTail, pSynced:
**   All dirty pages are linked into the doubly linked list using
**   PgHdr.pDirtyNext and pDirtyPrev. The list is maintained in LRU order
**   such that p was added to the list more recently than p.pDirtyNext.
**   PCache.pDirty points to the first (newest) element in the list and
**   pDirtyTail to the last (oldest).
*/
package bplustree

import (
  "unsafe"
  "C"
  //"fmt"
  //"encoding/binary"
)

/* Allowed values for second argument to ManageDirtyList() */
const (
  PCACHE_DIRTYLIST_REMOVE  = 1    /* Remove pgHdr from dirty list */
  PCACHE_DIRTYLIST_ADD     = 2    /* Add pgHdr to the dirty list */
  PCACHE_DIRTYLIST_FRONT   = 3    /* Move pgHdr to the front of the list */
  PGHDR_CLEAN = 1
  PGHDR_DIRTY = 2
)

type Bulk struct {
	addr uintptr
	len  int
	cap  int
}

type PageData Bulk
type CacheData Bulk

type PCache struct {
  szPage uint16                         /* Size of database content section */
  szAlloc uint16                     /* Total size of one pcache line */
  nMin uint32                  /* Minimum number of pages reserved */
  nMax uint32                  /* Configured "cache_size" value */
  pBulk *CacheData
  pLru *PgHdr
  pFree *PgHdr                     /* Next in hash table chain */
  pDirty *PgHdr                     /* Next in hash table chain */
  pDirtyTail *PgHdr

  /* Hash table of all pages. The following variables may only be accessed
  ** when the accessor is holding the PGroup mutex.
  */
  nPage uint32                 /* Total number of pages in apHash */
  nInitPage uint16
  nHash uint32                /* Number of slots in apHash[] */
  apHash []*PgHdr                    /* Hash table for fast lookup by key */
  iKey uint32                  /* Key value (pgHdr number) */
}

/*
** Every pgHdr in the cache is controlled by an instance of the following
** structure.
**
** A Page cache line looks like this:
**
**  --------------------------------------------------
**  |  database pgHdr content   |  PgHdr  |  PgHdr  |
**  --------------------------------------------------
*/
type PgHdr struct {
  iKey uint32                     /* Page number for this pgHdr */
  flag int                     /* Dirty of Clean*/
  isBulkLocal int
  pBulk *PageData                   /* Page data */

  pCache *PCache              /* PRIVATE: Cache that owns this pgHdr */
  pNext *PgHdr                 /* Transient list of dirty sorted by pgno */
  pFreeNext *PgHdr                 /* Transient list of dirty sorted by pgno */
  pDirtyNext *PgHdr             /* Next element in list of dirty pages */
  pDirtyPrev *PgHdr             /* Previous element in list of dirty pages */
  pLruNext *PgHdr             /* Next in LRU list of unpinned pages */
  pLruPrev *PgHdr              /* Previous in LRU list of unpinned pages */
}

/*
** Implementation of the Create method.
**
** Allocate a new cache.
*/
func (pCache *PCache) Create(szPage uint16) {
  pCache.nPage = 0

  pCache.szPage = szPage

  var pgHdr *PgHdr
  pCache.szAlloc = szPage + uint16(unsafe.Sizeof(*pgHdr))
  // pcache1EnterMutex(pGroup);
  pCache.ResizeHash()
  pCache.nMax = pCache.nHash
  // pcache1LeaveMutex(pGroup);
  if( pCache.nHash==0 ){
    pCache.Destroy()
  }
  pCache.InitBulk()
}

/*
** Try to initialize the pCache.pFree and pCache.pBulk fields.  Return
** true if pCache.pFree ends up containing one or more free pages.
*/
func (pCache *PCache) InitBulk() *PgHdr {
  /* Do not bother with a bulk allocation if the cache size very small */
  var szBulk int
  if pCache.nInitPage>0 {
    szBulk = int(pCache.szAlloc * pCache.nInitPage)
  } else {
    szBulk = int(pCache.szAlloc * 1024)
  }
  pBulk := C.malloc(C.size_t(szBulk))

  pCache.pBulk = &CacheData{
    addr: uintptr(unsafe.Pointer(pBulk)),
    len:  szBulk,
    cap:  szBulk,
  }

  nBulk := szBulk/int(pCache.szAlloc)
  for i:= uint16(0); i < uint16(nBulk); i++ {
    pX := (*PgHdr)(unsafe.Pointer(uintptr(unsafe.Pointer(pBulk)) +
    uintptr(i*pCache.szAlloc) + uintptr(pCache.szPage)))
    pX.pBulk = &PageData{
      addr: uintptr(unsafe.Pointer(pBulk))+uintptr(i*pCache.szAlloc),
      len:  int(pCache.szAlloc),
      cap:  int(pCache.szAlloc),
    }
    pX.pFreeNext = pCache.pFree
    pCache.pFree = pX

  }
  return pCache.pFree
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
  /* Step 2: If the pgHdr was found in the hash table, then return it.
  ** If the pgHdr was not in the hash table continue with
  ** subsequent steps to try to create the pgHdr. */
  pgHdr := pCache.apHash[iKey % pCache.nHash];
  for pgHdr != nil {
    if pgHdr.iKey == iKey {
      return pgHdr
    }
    pgHdr = pgHdr.pNext
  }

  /* Steps 3 if pgHdr num is nearly full resize the hash*/
  if pCache.nPage>=pCache.nHash {
    pCache.ResizeHash()
  }
  /* Step 4. Try to recycle a pgHdr. */
  if pCache.nPage+1 >= pCache.nMax  {
    pgHdr = pCache.pLru
    if pgHdr != nil {
      pCache.RemoveFromHash(pgHdr)
    }
  }
  /* Step 5. If a usable pgHdr buffer has still not been found,
  ** attempt to allocate a new one.
  */
  if pgHdr == nil {
    pgHdr = pCache.AllocPage()
  }

  if pgHdr != nil {
    h := iKey % pCache.nHash
    pCache.nPage++
    pgHdr.iKey = iKey
    pgHdr.pNext = pCache.apHash[h]

    pgHdr.pCache = pCache
    pgHdr.pLruPrev = nil
    pgHdr.pLruNext = pCache.pLru
    if pCache.pLru != nil {
      pCache.pLru.pLruPrev = pgHdr
    }
    pCache.pLru = pgHdr
    pCache.apHash[h] = pgHdr
  }
  return pgHdr;
}

/*
** Allocate a new pgHdr object initially associated with cache pCache.
*/
func (pCache *PCache) AllocPage() *PgHdr {
  if pCache.pFree != nil/*|| (pCache.nPage==0 && pcache1InitBulk(pCache))*/{
    pgHdr := pCache.pFree
    pCache.pFree = pgHdr.pFreeNext
    pgHdr.pFreeNext = nil
    return pgHdr
  }
  pBulk := C.malloc(C.size_t(pCache.szAlloc))
  pgHdr := (*PgHdr)(unsafe.Pointer(pBulk))
  pgHdr.pBulk = &PageData{
    addr: uintptr(unsafe.Pointer(pBulk)),
    len:  int(pCache.szAlloc),
    cap:  int(pCache.szAlloc),
  }
  pgHdr.isBulkLocal = 0
  return pgHdr
}

/*
** Free a pgHdr object allocated by pcache1AllocPage().
*/
func (pCache *PCache) FreePage(p *PgHdr){

  // if( p.isBulkLocal ){
  p.pFreeNext = pCache.pFree;
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

  for i:= uint32(0); i<pCache.nHash; i++{
    pCurPg := pCache.apHash[i];
    for pCurPg != nil {
      h := pCurPg.iKey % nNew;
      pNewPg := apNew[h]

      apNew[h] = pCurPg
      pCurPg = pCurPg.pNext
      apNew[h].pNext = pNewPg
    }
  }
  pCache.apHash = apNew;
  pCache.nHash = nNew;
}

/*
** Remove the pgHdr supplied as an argument from the hash table
** (PCache1.apHash structure) that it is currently stored in.
** Also free the pgHdr if freePage is true.
**
*/
func (pCache *PCache) RemoveFromHash(pgHdr *PgHdr) {

  h := pgHdr.iKey % pCache.nHash
  p := &pCache.apHash[h]
  for *p != nil {
    if (*p) == pgHdr {
      *p = (*p).pNext
      pCache.nPage--
      pCache.FreePage(pgHdr)
      return
    }
    p=&((*p).pNext)
  }
}

/*
** Manage pgHdr's participation on the dirty list.  Bits of the addRemove
** argument determines what operation to do.  The 0x01 bit means first
** remove pgHdr from the dirty list.  The 0x02 means add pgHdr back to
** the dirty list.  Doing both moves pgHdr to the front of the dirty list.
*/
func (pCache *PCache) ManageDirtyList(pgHdr *PgHdr, addRemove uint8){

  if addRemove & PCACHE_DIRTYLIST_REMOVE == 1 {

    /* Update the PCache.pSynced variable if necessary. */
    // if( p.pSynced==pgHdr ){
    //   p.pSynced = pgHdr.pDirtyPrev;
    // }

    if pgHdr.pDirtyNext != nil {
      pgHdr.pDirtyNext.pDirtyPrev = pgHdr.pDirtyPrev
    }else{
      pCache.pDirtyTail = pgHdr.pDirtyPrev
    }
    if pgHdr.pDirtyPrev != nil {
      pgHdr.pDirtyPrev.pDirtyNext = pgHdr.pDirtyNext
    }else{
      /* If there are now no dirty pages in the cache, set eCreate to 2.
      ** This is an optimization that allows sqlite3PcacheFetch() to skip
      ** searching for a dirty pgHdr to eject from the cache when it might
      ** otherwise have to.  */
      pCache.pDirty = pgHdr.pDirtyNext
    }
    pgHdr.pDirtyNext = nil
    pgHdr.pDirtyPrev = nil
  }
  if addRemove & PCACHE_DIRTYLIST_ADD == 1 {
    pgHdr.pDirtyNext = pCache.pDirty;
    if pgHdr.pDirtyNext != nil {
      pgHdr.pDirtyNext.pDirtyPrev = pgHdr;
    }else{
      pCache.pDirtyTail = pgHdr;
    }
    pCache.pDirty = pgHdr;
  }
}

/*
** Make sure the pgHdr is marked as dirty. If it isn't dirty already,
** make it so.
*/
func (pCache *PCache) MakeDirty(pgHdr *PgHdr){
  if pgHdr.flag & PGHDR_CLEAN != 0 {
    pgHdr.flag ^= (PGHDR_DIRTY|PGHDR_CLEAN)
    pCache.ManageDirtyList(pgHdr, PCACHE_DIRTYLIST_ADD)
  }
}

/*
** Make sure the pgHdr is marked as clean. If it isn't clean already,
** make it so.
*/
func (pCache *PCache) MakeClean(pgHdr *PgHdr){
  if (pgHdr.flag & PGHDR_DIRTY) != 0 {
    pCache.ManageDirtyList(pgHdr, PCACHE_DIRTYLIST_REMOVE)
    pgHdr.flag &= ^(PGHDR_DIRTY)
    pgHdr.flag |= PGHDR_CLEAN
  }
}

/*
** Make every pgHdr in the cache clean.
*/
func (pCache *PCache) MakeCleanAll(){
  for pCache.pDirty != nil {
    p := pCache.pDirty
    pCache.MakeClean(p)
  }
}
