package crawler

import (
  "fmt"
  "github.com/opesun/goquery"
)

type crawler struct {
  runStatus bool
  urlChan chan []string
}

func (c *crawler)doCrawl(url string) {
  // login in weibo

  for c.runStatus {
    select {
    case x := urlChan
      //do something
    default:
    //do something
    }
  }
}
