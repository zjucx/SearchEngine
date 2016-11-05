package segment

import (
  "fmt"
  "os"
  "path/filepath"
  "bufio"
  "io"
  "strings"
  "github.com/yanyiwu/gojieba"
  "regexp"
  "invertidx"
)

func Segment(){
  segment("./tmp")
}

func segment(path string) {
  //build dict for index
  dict := invertidx.NewDict("./src/invertidx/tmpfile/dict.dict")
  stopdict := invertidx.NewDict("./src/invertidx/tmpfile/stop.dict")
  idxBuf := invertidx.NewIdxBuf("./src/invertidx/tmpfile/idx.tmp")

  //walk the file in specific dir and write to the raw index file
  err := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
    if ( fi == nil ) {return err}
    if fi.IsDir() {return nil}
    f, err := os.Open(path)
    if err != nil {return err}
    buf := bufio.NewReader(f)
    for {
      line, err := buf.ReadString('\n')
      //去除非中文字符
      reg := regexp.MustCompile("[^\u4e00-\u9fa5]")
      line = reg.ReplaceAllString(line, "")
      //jieba分词
      var words []string
      use_hmm := true
      x := gojieba.NewJieba()
      defer x.Free()
      words = x.Cut(line, use_hmm)
      //fmt.Println(line)
      fmt.Println("精确模式:", strings.Join(words, "/"))

      for _, v := range words {
        //build dictory
        v := dict.AddDict(v)
        fmt.Println(v)
        idxBuf.AddIndexItem(1, v)
      }
      //segLine(line, dict, idxBuf)

      if err != nil {
  			if err == io.EOF {
  				return nil
  			}
  			return err
  		}
    }
  })
  if err != nil {
    fmt.Printf("filepath.Walk() returned %v\n", err)
  }
}

/*func segLine(line string, dict, idxBuf interface{}) {
  var words []string
  use_hmm := true
  x := gojieba.NewJieba()
  defer x.Free()
  words = x.Cut(line, use_hmm)
  //fmt.Println(line)
  fmt.Println("精确模式:", strings.Join(words, "/"))

  for _, v := range words {
    //build dictory
    v := dict.AddDict(v)
    fmt.Println(v)
    idxBuf.AddIndexItem(1, v)
  }
}*/
