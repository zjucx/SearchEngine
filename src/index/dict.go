package index

import (
  "fmt"
  "os"
  "bufio"
)

type Dictionary struct {
  dict map[string]int
  curSize int
  filename string
}

func (d *dictionary)AddDict(key string) int {
  strkey := []byte(key)
  if len(key) > maxkeylen {
    strkey = strings(key[0:maxkeylen-1])
  }

  // 查找键值是否存在
  if v, ok := d.dict[key]; !ok {
    d.dick[key] = ++d.curSize
  	//fmt.Println(v)
    return d.curSize
  }
  return d.dict[key]
}

func (d *dictionary)LoadDictFile() error {
  fi, err := os.Open(d.filename)
  if err != nil {
    panic(err)
  }
  defer fi.Close()
  //dict := make(map[string]uint32)
  br = bufio.NewReader(fi)
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
    d.dict[words[0]] = words[1]
	}
  d.curSize = len(d.dict)
  return nil
}

func (d *dictionary)WriteDictFile() error{
  if checkFileIsExist(d.filename) {  //如果文件存在
    fo, _ = os.OpenFile(d.filename, os.O_APPEND, 0666)  //打开文件
    fmt.Println("file exist!");
   }else {
    fo, _ = os.Create(d.filename)  //创建文件
    fmt.Println("file is not exist!");
   }
   /*strkey := []byte(key)
   if len(key) > maxkeylen {
     strkey = strings(key[0:maxkeylen-1])
   }*/
   bw = bufio.NewWriter(fo)
   buf := make([]byte, maxkeylen+6)
   for k, v := range m1 {
     strkey = strings(k[0:maxkeylen])
     strval := strconv.Itoa(v)
     str := strkey + "," + strval + "\n"//strings.Join({strkey, strval}, ",")
     err := bw.WriteFile(str)  //写入文件(字节数组)
     if !err {
       fmt.Println("write file error!");
       return err
     }
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
