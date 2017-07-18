package bplustree

import (
  "unsafe"
)
/*
** An instance of this object represents a single database file.
**
** A single database file can be in use at the same time by two
** or more database connections.  When two or more connections are
** sharing the same database file, each connection has it own
** private Btree object for the file and each of those Btrees points
** to this one BPlusTree object.
*/
type BPlusTree struct {
  Pager *pPager        /* The page cache */
  MemPage *pPage      /* First page of the database */
  uint16 maxLocal         /* Maximum local payload in non-LEAFDATA tables */
  uint16 minLocal         /* Minimum local payload in non-LEAFDATA tables */
  uint16 maxLeaf          /* Maximum local payload in a LEAFDATA table */
  uint16 minLeaf          /* Minimum local payload in a LEAFDATA table */
  uint32 pageSize         /* Total number of bytes on a page */
  uint32 usableSize       /* Number of usable bytes on each page */
  uint32 nPage            /* Number of pages in the database */
  num2page map[uint32]*MemPage
}

/* Each btree pages is divided into three sections:  The header, the
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
**      0       1      Flags. 1: interpage, 2: leafpage, 4: overflowpage, 8: zeropage
**      1       2      byte offset to the first freeblock
**      3       2      number of cells on this page
**      5       2      first byte of the cell content area
**      7       1      number of fragmented free bytes
**      8       4      Right child (the Ptr(N) value).  Omitted on leaves.
*/
type PageHeader struct {
  flag uint8
  freeOffset uint16
  nCell uint16
  plOffset uint16
  nFree uint8
  parent uint32
}

type MemPage struct{
  uint16 pgno           /* Page number for this page */
  /* Only the first 8 bytes (above) are zeroed by pager.c when a new page
  ** is allocated. All fields that follow must be initialized before use */
  u8 leaf             /* True if a leaf page */
  u8 hdrOffset        /* 100 for page 1.  0 otherwise */
  uint16 nFree           /* Number of free bytes on the page */
  uint16 nCell           /* Number of cells on this page, local and ovfl */
  u8 *aData           /* Pointer to disk image of the page data */
  u8 *aDataEnd        /* One byte past the end of usable data */
  u8 *aCellIdx        /* The cell index area */
  u8 *aDataOfst       /* Same as aData for leaves.  aData+4 for interior */
}


/* The basic idea is that each page of the file contains N database
** entries and N+1 pointers to subpages.
**
**   ----------------------------------------------------------------
**   |  Ptr(0) | Key(0) | Ptr(1) | Key(1) | ... | Key(N-1) | Ptr(N) |
**   ----------------------------------------------------------------
*/
type Cell struct {
  ptr      uint32      /* page number */
  key      uint32      /* The key for Payload or Offset of the page start of a payload*/
}

/* DocId1 DocId2 ...
**   -----------------------------------------------------------------
**   |  key | DocId1 | DocId1 | DocId3 | ... | DocId(N-1) | DocId(N) |
**   -----------------------------------------------------------------
 */
type Payload struct {
  key     uint32             /* value in the unpacked key */
  entrys  unsafe.Pointer            /* fot data compress */
  nEntry  uint16             /* Number of values.  Might be zero */
}
