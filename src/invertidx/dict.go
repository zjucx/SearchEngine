package invertidx

import (
  "os"
  "bufio"
  "io"
  "strings"
  "strconv"
  // "fmt"
)

type Dictionary struct {
  dict map[string]int
  curSize int
  filename string
}

const maxkeylen = 64

func NewDict(filename string) Dictionary{
  dict := Dictionary{
    filename : filename,
    dict : make(map[string]int),
  }
  if checkFileIsExist(filename) {
    dict.LoadDictFile()
  }

  return dict
}

func (d *Dictionary)AddDict(key string) int {
  if len(key) > maxkeylen {
    key = string(key[0:maxkeylen-1])
  }
  // 查找键值是否存在
  if v, ok := d.dict[key]; ok {
    return v
  } else {
    d.curSize++
    d.dict[key] = d.curSize
    //fmt.Println(v)
    return d.curSize
  }
}

func (d *Dictionary)LoadDictFile() error {
  fi, err := os.Open(d.filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  //dict := make(map[string]uint32)
  br := bufio.NewReader(fi)
  //read buffer to map
  for {
		line, err := br.ReadString('\n')
		//line = strings.TrimSpace(line)
		//handler(line)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

    words := strings.Split(line, ",")
    wordid, _ := strconv.Atoi(words[1])
    d.dict[words[0]] = wordid
	}
  d.curSize = len(d.dict)
  return nil
}

func (d *Dictionary)WriteDictFile(){

   fo, err := OpenFile(d.filename)
   CheckErr(err)
   /*strkey := []byte(key)
   if len(key) > maxkeylen {
     strkey = strings(key[0:maxkeylen-1])
   }*/
   bw := bufio.NewWriter(fo)
   //buf := make([]byte, maxkeylen+6)
   for k, v := range d.dict {
     strval := strconv.Itoa(v)
     str := k + "," + strval + "\n"//strings.Join({strkey, strval}, ",")
     _, err := bw.WriteString(str)  //写入文件(字节数组)
     CheckErr(err)
   }

}

/**
 * 判断文件是否存在  存在返回 true 不存在返回false
 */
func checkFileIsExist(filename string) (bool) {
 var exist = true;
 if _, err := os.Stat(filename); os.IsNotExist(err) {
  exist = false;
 }
 return exist;
}
