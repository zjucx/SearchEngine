package scrawler

import (
  "fmt"
  "regexp"
  "strconv"
  "strings"
  "os"
  "bufio"
)

var mstartUrl = "http://d.weibo.com/1087030002_2975_1003_0"

func Scrawler(username, passwd string){
  // get login cookies
  loginCookies := WeiboLogin(username, passwd)
  mUrl := make(map[string]string)
  for i := 1; i < 168; i++ {
    mstartUrl = "http://d.weibo.com/1087030002_2975_1003_0?pids=Pl_Core_F4RightUserList__4&page=" + strconv.Itoa(i) + "#Pl_Core_F4RightUserList__4"
    mstartResp, _ := DoRequest(`GET`, mstartUrl, ``, loginCookies, ``, header)

    // get

    //fmt.Println(mstartResp)
    //resp make '\' as string so we shou add '\\' in regex and one blank become two blanks so be carefully
    //the regex should reg  the string of mstartResp
    reg := regexp.MustCompile(`<a class=\\"S_txt1\\" target=\\"_blank\\"  usercard=\\"(.*?)\\" href=\\"(.*?)\\" title=\\"(.*?)\\">`)
    arrStart := reg.FindAllStringSubmatch(mstartResp, -1)

    if len(arrStart) > 0 {
      fmt.Println(len(arrStart))
      for i := 0; i < len(arrStart); i++ {
        mUrl[arrStart[i][3]] = strings.Replace(arrStart[i][2], "\\/","/", -1)
        fmt.Println(strings.Replace(arrStart[i][2], "\\/","/", -1) + arrStart[i][3])
      }
    }
  }
  writeMaptoFile(mUrl, "./mstarturlname.map")
}

func writeMaptoFile(m map[string]string, filePath string) error {
  f, err := os.Create(filePath)
  if err != nil {
          fmt.Printf("create map file error: %v\n", err)
          return err
  }
  defer f.Close()

  w := bufio.NewWriter(f)
  for k, v := range m {
          lineStr := fmt.Sprintf("%s^%s", k, v)
          fmt.Fprintln(w, lineStr)
  }
  return w.Flush()
}
