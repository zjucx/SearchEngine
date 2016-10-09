package page

import (
  "fmt"
  "unsafe"
)

const maxkeylen uint32 = 64

const (
	tindex     uint8 = 0
	tdata      uint8 = 1
)

type item struct {
  k [maxkeylen]byte  //key for
  count uint32       //count of docid for this key
  offset uint32      //the offset of key in data file offset+sizeof(uint32)
}

type page struct {
  pid uint32
  ptype uint8       //dacid data page and index page
  data uintptr      //save data for docid or index page
}

func (p *page) getDataPtr {
  return p.ptype ? (*uint32)(unsafe.Pointer(&p.data)) : &((*[0xFFFF]item)(unsafe.Pointer(&p.data)))[0]
}
