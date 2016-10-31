package invertidx

import (
  //"fmt"
  "os"
)
func OpenFile(filename string) (*os.File, error) {
  f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
  return f, err
}

func CheckErr(err error) {
  if err != nil {
    panic(err)
  }
}
