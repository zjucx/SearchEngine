package bplustree

import (
  "fmt"
)

// Each btree pages is divided into three sections:  The header, the
// cell pointer array, and the cell content area.  Page 1 also has a 100-byte
// file header that occurs before the page header.
//
//      |----------------|
//      | file header    |   100 bytes.  Page 1 only.
//      |----------------|
//      | page header    |   8 bytes for leaves.  12 bytes for interior nodes
//      |----------------|
//      | cell pointer   |   |  2 bytes per cell.  Sorted order.
//      | array          |   |  Grows downward
//      |                |   v
//      |----------------|
//      | unallocated    |
//      | space          |
//      |----------------|   ^  Grows upwards
//      | cell content   |   |  Arbitrary order interspersed with freeblocks.
//      | area           |   |  and free space fragments.
//      |----------------|
//
// The page headers looks like this:
//
//   OFFSET   SIZE     DESCRIPTION
//      0       1      Flags. 1: intkey, 2: zerodata, 4: leafdata, 8: leaf
//      1       2      byte offset to the first freeblock
//      3       2      number of cells on this page
//      5       2      first byte of the cell content area
//      7       1      number of fragmented free bytes
//      8       4      Right child (the Ptr(N) value).  Omitted on leaves.
//

type MemPage struct{
  u8 intKey;           /* True if table b-trees.  False for index b-trees */
  u8 intKeyLeaf;       /* True if the leaf of an intKey table */
  Pgno pgno;           /* Page number for this page */
  /* Only the first 8 bytes (above) are zeroed by pager.c when a new page
  ** is allocated. All fields that follow must be initialized before use */
  u8 leaf;             /* True if a leaf page */
  u8 hdrOffset;        /* 100 for page 1.  0 otherwise */
  u8 childPtrSize;     /* 0 if leaf==1.  4 if leaf==0 */
  u8 max1bytePayload;  /* min(maxLocal,127) */
  u8 nOverflow;        /* Number of overflow cell bodies in aCell[] */
  u16 maxLocal;        /* Copy of BtShared.maxLocal or BtShared.maxLeaf */
  u16 minLocal;        /* Copy of BtShared.minLocal or BtShared.minLeaf */
  u16 cellOffset;      /* Index in aData of first cell pointer */
  u16 nFree;           /* Number of free bytes on the page */
  u16 nCell;           /* Number of cells on this page, local and ovfl */
  u16 maskPage;        /* Mask for page offset */
  u16 aiOvfl[4];       /* Insert the i-th overflow cell before the aiOvfl-th
                       ** non-overflow cell */
  u8 *apOvfl[4];       /* Pointers to the body of overflow cells */
  u8 *aData;           /* Pointer to disk image of the page data */
  u8 *aDataEnd;        /* One byte past the end of usable data */
  u8 *aCellIdx;        /* The cell index area */
  u8 *aDataOfst;       /* Same as aData for leaves.  aData+4 for interior */
  DbPage *pDbPage;     /* Pager page handle */
  u16 (*xCellSize)(MemPage*,u8*);             /* cellSizePtr method */
  void (*xParseCell)(MemPage*,u8*,CellInfo*); /* btreeParseCell method */
}

type Payload struct {
  key    int               /* PRIMARY KEY for tabs */
  values []int             /* value in the unpacked key */
  nValue   int             /* Number of values.  Might be zero */
}

type CellInfo struct {
  i64 nKey;      /* The key for INTKEY tables, or nPayload otherwise */
  u8 *pPayload;  /* Pointer to the start of payload */
  u32 nPayload;  /* Bytes of payload */
  u16 nLocal;    /* Amount of payload held locally, not on overflow */
  u16 nSize;     /* Size of the cell content on the main b-tree page */
};

type BTree struct {
  root *InterNode
  frist *LeafNode
  leafCount int
  interCount int
  height int
}

func NewBTree() *BTree {
  leaf := NewLeafNode(nil)
  r := NewInterNode(nil, leaf)
  leaf.parent = r

  return &BTree{
    root : r,
    frist : leaf,
    leafCount : 1,
    interCount : 1,
    height : 2,
  }
}


func (bt *BTree) Insert(key int, value string) {
  _, index, l := search(bt.root, key)
  leaf, k, ok := l.Insert(key, value)
  if !ok {
    return
  }
  var newNode node
  newNode = leaf
  p := l.getParent()
  for {
    n, k, ok := p.Insert(k, newNode)
    if !ok {
      return
    }
    isRoot := p.parent == nil
    if isRoot {
      r := newInterNode(k, n)
      r.Insert(p.kcs[0].key, p)
      return
    }
    p = p.getParent()
    newNode = n
  }
}

func (bt *BTree) Search(key int) (string, bool) {
  kv, _, _ := search(bt.root, key)
  if kv == nil {
    return "", false
  }
  return kv.value, true
}

func search(bt *BTree, key int) (*kv, int, *LeafNode) {
  curr := bt.root
  index := -1
  for {
    switch t := curr.(type) {
    case *LeafNode:
      i, ok := t.find(key)
      if !ok {
        return "", index, t
      }
      return t.kvs[i].value, index, t
    case *InterNode:
      i, _ = t.find(key)
      index = i
      curr = t.kcs[i].child
    default:
      panic("")
    }
  }
}
