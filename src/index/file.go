package index

import (
  "fmt"
  "os"
  "bufio"
)
func OpenFile(filename string) (*File) {
  f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
  if err != nil {
      b.Error("Open file error")
      return nil
  }
  //defer f.Close()
  return f
}

func CheckErr(err error) {
  if err != nil {
    panic(err)
  }
}
