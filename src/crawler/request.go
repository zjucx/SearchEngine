package scrawler

import (
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
)
/**
 *
 */
func DoRequest(method, reqUrl, params, cookies, domain string, header map[string]string) (resBody, resCookies string){
  //prepare params for method
  method = strings.ToUpper(method)

  //prepare params for reader
  var paramsReader io.Reader
  if method == "POST" {
    paramsReader = strings.NewReader(params)
  } else {
    paramsReader = nil
  }

  // 1) complete request
  req, err := http.NewRequest(method, reqUrl, paramsReader)
  if err != nil {
    log.Fatalln("NewRequest Err:", err)
    return
  }

  // 2) prepare header if has
  if header != nil {
    for k, v := range header {
      req.Header.Set(k, v)
    }
  }

  // 3) prepare cookie if has
  gCookieJar, _ := cookiejar.New(nil)
  if len(cookies) > 0 {
    cookies := appendCookies(cookies, "", domain)
		cookieUrl, _ := url.Parse(reqUrl)
		gCookieJar.SetCookies(cookieUrl, cookies)
  }

  // 4) prepare transport
  var transport *http.Transport = nil
  //proxyUrl, _ := url.Parse("http://127.0.0.1:8888")
  if true {
    transport = &http.Transport{
      TLSClientConfig: &tls.Config{
        InsecureSkipVerify: true,
        MinVersion:         tls.VersionTLS10,
        MaxVersion:         tls.VersionTLS12,
      },
      DisableCompression:true,
      //Proxy: http.ProxyURL(proxyUrl),
    }
  } else {
    transport = &http.Transport{}
  }

  // construct req client
  client := &http.Client{
		Transport:     transport,
		CheckRedirect: nil,
		Jar:           gCookieJar,
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln("Do Request Err:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("ReadAll Response Err:", err)
		return
	}
	resBody = string(body)

	arrCookies := resp.Cookies()

	for _, data := range arrCookies {
		resCookies += data.Name + "=" + data.Value + ";"
	}

	/*if len(resCookies) > 0 {
		resCookies = SubString(respCookies, 0, len(respCookies)-1)
	}*/

	return
}

/*
 * @functional http请求附加cookie
 * @param string strCookies
 * @return []*http.Cookie
 */
func appendCookies(strCookies, path, domain string) []*http.Cookie {
	var cookies []*http.Cookie

	if path == "" {
		path = "/"
	}

  // parse cookie from string to map
  mapCookie := make(map[string]string)
	reg := regexp.MustCompile(`([^=]+)=([^;]*);?`)
	arrCookie := reg.FindAllStringSubmatch(strCookies, -1)
	if len(arrCookie) > 0 {
		for i := 0; i < len(arrCookie); i++ {
			mapCookie[arrCookie[i][1]] = arrCookie[i][2]
		}
	}

	for k, v := range mapCookie {
		appendCookie := &http.Cookie{
			Name:   k,
			Value:  v,
			Path:   path,
			Domain: domain,
		}
		cookies = append(cookies, appendCookie)
	}
	return cookies
}
