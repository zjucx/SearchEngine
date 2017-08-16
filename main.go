package main

import (
	//"fmt"
	//"os"
	//"segment"
	//"invertidx"
	//"web"
	"bplustree"
	//"unsafe"
	//"C"
	//"fmt"
)
type SliceHeader struct {
	addr uintptr
	len  int
	cap  int
}

func main() {
		//scrawler.Scrawler("xxxxxxx@163.com", "xxxxxxx")
		//segment.Segment()
		//scrawler.Scrawler()
		//web.Main()
		bpTree := &bplustree.BPlusTree{}
		bpTree.Open("test.db", 1024)

		/*pBulk := C.malloc(C.size_t(1024))//make([]byte, szBulk)

		sl := &SliceHeader{
		  addr: uintptr(unsafe.Pointer(pBulk)),
		  len:  1024,
		  cap:  1024,
		 }
	  nBulk := 4
		p := *(*[]byte)(unsafe.Pointer(sl))
	  for i:= 0; i < nBulk; i++ {
	    fmt.Println(p[1024])
			println(len(p))
	  }*/
}
