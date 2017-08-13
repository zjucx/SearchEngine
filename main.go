package main

import (
	//"fmt"
	//"os"
	//"segment"
	//"invertidx"
	//"web"
	"bplustree"
)

func main() {
		//scrawler.Scrawler("xxxxxxx@163.com", "xxxxxxx")
		//segment.Segment()
		//scrawler.Scrawler()
		//web.Main()
		bpTree := &BPlusTree{}
		bpTree.Open("test.db", 1024)
}
