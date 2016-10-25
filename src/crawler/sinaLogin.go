package scrawler

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/url"
	"os"
  "regexp"
  "encoding/json"
  "bufio"
  "strconv"
)

var header = map[string]string{
  "Host":                      "login.sina.com.cn",
  "Proxy-Connection":          "keep-alive",
  "Cache-Control":             "max-age=0",
  "Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
  "Origin":                    "http://weibo.com",
  "Upgrade-Insecure-Requests": "1",
  "User-Agent":                "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.94 Safari/537.36",
  "Referer":                   "http://weibo.com",
  "Accept-Language":           "zh-CN,zh;q=0.8,en;q=0.6,ja;q=0.4",
  "Content-Type":              "application/x-www-form-urlencoded",
}

func WeiboLogin(username, passwd string){
  //get cookie for sina website
  strCookies := getCookies()
  // crypto username for logining
  su := url.QueryEscape(username)
  su = base64.StdEncoding.EncodeToString([]byte(su))

  // crypto password for logining
  loginInfo := getPreLogin(su)
  sp := encryptPassword(loginInfo, passwd)

  // is need cgi or not
  var cgi string
  if loginInfo["showpin"] == 1 {
    inputDone := make (chan string)
    go inputcgi(inputDone)
    cgi = <- inputDone
  }

  // Do login POST
  loginUrl := `http://login.sina.com.cn/sso/login.php?client=ssologin.js(v1.4.18)`
  // form data params
  strParams := buildParems(su, sp, cgi, loginInfo)
  loginResp, loginCookies := DoRequest(`POST`, loginUrl, strParams, strCookies, ``, header)
  //请求passport
	passportResp, _ := callPassport(loginResp, strCookies+";"+loginCookies)
	uniqueid := MatchData(passportResp, `"uniqueid":"(.*?)"`)
	homeUrl := "http://weibo.com/u/" + uniqueid + "/home?topnav=1&wvr=6"

	//进入个人主页
	entryHome(homeUrl, loginCookies)
	//抓取个首页
	//result := getPage(loginCookies)
	//fmt.Println(result)
}

func inputcgi(inputDone chan string){
  reader := bufio.NewReader(os.Stdin)
  //for {
  fmt.Println("waiting for input captcha...")
  data, _, _ := reader.ReadLine()
  inputDone <- string(data)
  //}
}

/*
 * crypto passwd for logining
 * var RSAKey = new sinaSSOEncoder.RSAKey();
 * RSAKey.setPublic(me.rsaPubkey, "10001");
 * password = RSAKey.encrypt([me.servertime, me.nonce].join("\t") + "\n" + password)
 *
 */
func encryptPassword(loginInfo map[string]interface{}, password string) string {
  fmt.Println("======encryptPassword")
  z := new(big.Int)
	z.SetString(loginInfo["pubkey"].(string), 16)
	pub := rsa.PublicKey{
		N: z,
		E: 65537,
	}
	encryString := strconv.Itoa(int(loginInfo["servertime"].(float64))) + "\t" + loginInfo["nonce"].(string) + "\n" + password
	encryResult, _ := rsa.EncryptPKCS1v15(rand.Reader, &pub, []byte(encryString))
	return hex.EncodeToString(encryResult)
}

/*
 * open main page and you should get cookie and save
 */
 func getCookies() string{
   fmt.Println("======getCookies")
   loginUrl := `http://weibo.com/login.php`
   _, strCookies := DoRequest(`GET`, loginUrl, ``, ``, ``, nil)
   return strCookies
 }

/*
 * when finish inputing the username, send the prelogin req
 * you can get login info for logining sina
 */
func getPreLogin(su string) map[string]interface{} {
  preLoginUrl := `https://login.sina.com.cn/sso/prelogin.php?entry=weibo&callback=sinaSSOController.preloginCallBack&su=`+
  su + `&rsakt=mod&checkpin=1&client=ssologin.js(v1.4.18)&_=`
  resBody, _ := DoRequest(`GET`, preLoginUrl, ``, ``, ``, nil)
  //use regex extra json string
  strLoginInfo := RegexFind(resBody, `\((.*?)\)`)
  fmt.Println("======getPreLogin:" + strLoginInfo)
  //parse json str to map[string]string
  //json str 转map
	var loginInfo map[string]interface{}
	if err := json.Unmarshal([]byte(strLoginInfo), &loginInfo); err == nil {
		fmt.Println("==============json str 转map=======================")
		fmt.Println(loginInfo["pubkey"].(string))
    //return nil
	}
  return loginInfo
}

/*
 * entry:weibo
 * gateway:1
 * from:
 * savestate:7
 * useticket:1
 * pagerefer:
 * vsnf:1
 * su:aGZ1dGN4JTQwMTYzLmNvbQ==
 * service:miniblog
 * servertime:1477206529
 * nonce:2D9O10
 * pwencode:rsa2
 * rsakv:1330428213
 * sp:b96481646e643b59373c8b706e439c5f5b95990b7110e62e7f7e67ccab81571fc2e216950c6bf5764e181c2735839eb161d074ea489d2254be4a6756e05745a5fde469f30d3ae23539d1c74d321f08fc169e08f2f5da9f49c9f7e40e17c5a3d278b6bfcca214c70ed4fd37cb75c8d0e4a8d30fe671c418fc5a256305c93bafd0
 * sr:1280*800
 * encoding:UTF-8
 * prelt:839
 * url:http://weibo.com/ajaxlogin.php?framelogin=1&callback=parent.sinaSSOController.feedBackUrlCallBack
 * returntype:META
 */
func buildParems(su, sp, captcha string, loginInfo map[string]interface{}) string {
  fmt.Println("======buildParems")
  strParams := "entry=weibo&gateway=1&from=&savestate=7&useticket=1&pagerefer=&vsnf=1&su=" +
  su + "&service=miniblog&servertime=" + strconv.Itoa(int(loginInfo["servertime"].(float64))) +
  "&nonce=" + loginInfo["nonce"].(string) +
  "&pwencode=rsa2&rsakv=" + loginInfo["rsakv"].(string) +
  "&sp=" + sp +
  "&sr=1280*800&encoding=UTF-8&prelt=839&url=http%3A%2F%2Fweibo.com%2Fajaxlogin.php%3Fframelogin%3D1%26callback%3Dparent.sinaSSOController.feedBackUrlCallBack&returntype=META"
  //需要验证码
	if loginInfo["showpin"].(float64) == 1 {
		strParams += "&door=" + captcha
	}
  return strParams
}

//获取passport并请求
func callPassport(resp, cookies string) (passresp, passcookies string) {
  fmt.Println("======callPassport:" + resp)
	//提取passport跳转地址
	passportUrl := RegexFind(resp, `location.replace\(\'(.*?)\'\)`)
  fmt.Println("======callPassport:" + passportUrl)
	passresp, passcookies = DoRequest(`GET`, passportUrl, ``, cookies, ``, header)
	return
}

//进入首页
func entryHome(redirectUrl, cookies string) (homeResp, homeCookies string) {
  fmt.Println("======entryHome" + redirectUrl + cookies)
	homeResp, homeCookies = DoRequest(`GET`, redirectUrl, ``, cookies, ``, header)
	return
}

/*
 * @functional 正则表达式提取数据
 * @param string strText 输入文本
 * @param string strReg 正则表达式
 * @return string
 */
func RegexFind(strText, strReg string) (result string) {
	reg := regexp.MustCompile(strReg)
	arrMatch := reg.FindAllStringSubmatch(strText, -1)
	if len(arrMatch) > 0 {
		result = arrMatch[0][1]
	}
	return
}
/**
 * @functional 正则表达式匹配数据
 * @string strText 源字符串
 * @string strReg 正则表达式
 * @return string
 */
func MatchData(strText, strReg string) (result string) {
	reg := regexp.MustCompile(strReg)
	arrMatch := reg.FindAllStringSubmatch(strText, -1)
	if len(arrMatch) > 0 {
		result = arrMatch[0][1]
	}
	return
}
