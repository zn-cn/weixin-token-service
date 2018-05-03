package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type weixinAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Errcode     int    `json:"errcode"`
	Errmsg      string `json:"errmsg"`
}

var appid string
var tokenMutex sync.RWMutex
var token = ""
var tokenCH chan string
var jsapiTicketMutex sync.RWMutex
var jsapiTicket = ""
var accessTokenURL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
var jsapiTicketURL = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes 生成随机字符串
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getAccessTokenURL(appid, secret string) string {
	return fmt.Sprintf(accessTokenURL, appid, secret)
}

// WXConfigSign 生成签名
func WXConfigSign(jsapiTicket, nonceStr, timestamp, url string) (signature string) {
	if i := strings.IndexByte(url, '#'); i >= 0 {
		url = url[:i]
	}

	n := len("jsapi_ticket=") + len(jsapiTicket) +
		len("&noncestr=") + len(nonceStr) +
		len("&timestamp=") + len(timestamp) +
		len("&url=") + len(url)
	buf := make([]byte, 0, n)

	buf = append(buf, "jsapi_ticket="...)
	buf = append(buf, jsapiTicket...)
	buf = append(buf, "&noncestr="...)
	buf = append(buf, nonceStr...)
	buf = append(buf, "&timestamp="...)
	buf = append(buf, timestamp...)
	buf = append(buf, "&url="...)
	buf = append(buf, url...)

	hashsum := sha1.Sum(buf)
	return hex.EncodeToString(hashsum[:])
}

// views
func pingView(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func accessTokenView(c echo.Context) error {
	tokenMutex.RLock()
	defer tokenMutex.RUnlock()
	return c.JSON(http.StatusOK, map[string]interface{}{
		"access_token": token,
	})
}

func requestAccessToken(appid, secret string) (*weixinAccessTokenResponse, error) {
	resp, err := http.Get(getAccessTokenURL(appid, secret))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respBody weixinAccessTokenResponse
	err = json.Unmarshal(content, &respBody)
	if respBody.Errcode != 0 {
		return &respBody, fmt.Errorf("%d: %s", respBody.Errcode, respBody.Errmsg)
	}
	return &respBody, err
}

func cacheAccessToken(appid, secret string) (string, int) {
	log.Println("[wxtoken.cacheAccessToken] caching:main requesting new token")
	resp, err := requestAccessToken(appid, secret)
	if err != nil {
		log.Printf("[wxtoken.cacheAccessToken] caching:main request failed %v\n", err)
		return "", 1
	}
	tokenMutex.Lock()
	defer tokenMutex.Unlock()
	token = resp.AccessToken
	log.Printf("[wxtoken.cacheAccessToken] caching:main get response %v\n", resp)
	return resp.AccessToken, resp.ExpiresIn
}

// Config 全局配置结构体
type Config struct {
	AppID     string
	AppSecret string
	Addr      string
}

func getConfig() *Config {
	c := &Config{
		AppID:     "",
		AppSecret: "",
		Addr:      ":3001",
	}

	if v, ok := os.LookupEnv("WXTOKEN_APPID"); ok {
		c.AppID = v
	}

	if v, ok := os.LookupEnv("WXTOKEN_APPSECRET"); ok {
		c.AppSecret = v
	}

	if v, ok := os.LookupEnv("WXTOKEN_ADDR"); ok {
		c.Addr = v
	}

	return c
}

type weixinJSApiTicketResponse struct {
	Ticket    string `json:"ticket"`
	ExpiresIn int    `json:"expires_in"`
	Errcode   int    `json:"errcode"`
	Errmsg    string `json:"errmsg"`
}

func jsapiTicketView(c echo.Context) error {
	jsapiTicketMutex.RLock()
	defer jsapiTicketMutex.RUnlock()
	return c.JSON(http.StatusOK, map[string]string{
		"jsapi_ticket": jsapiTicket,
	})
}

func jssdkConfigView(c echo.Context) error {
	url := c.QueryParam("url")

	nonceStr := RandStringRunes(32)
	timestamp := time.Now().Unix()
	timestampStr := fmt.Sprintf("%d", timestamp)

	jsapiTicketMutex.RLock()
	defer jsapiTicketMutex.RUnlock()
	sign := WXConfigSign(jsapiTicket, nonceStr, timestampStr, url)

	return c.JSON(http.StatusOK, map[string]string{
		"appId":     appid,
		"nonce_str": nonceStr,
		"signature": sign,
		"timestamp": timestampStr,
	})
}

func getJSApiTicketURL(accessToken string) string {
	return fmt.Sprintf(jsapiTicketURL, accessToken)
}

func requestJSApiTicket(accessToken string) (*weixinJSApiTicketResponse, error) {
	resp, err := http.Get(getJSApiTicketURL(accessToken))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var respBody weixinJSApiTicketResponse
	err = json.Unmarshal(content, &respBody)
	if respBody.Errcode != 0 {
		return &respBody, fmt.Errorf("%d: %s", respBody.Errcode, respBody.Errmsg)
	}
	return &respBody, err
}

func cacheJSApiTicket(accessToken string) (string, int) {
	log.Println("requesting new jsapi ticket")
	resp, err := requestJSApiTicket(accessToken)
	if err != nil {
		log.Println("[wxtoken.cacheJSApiTicket] caching:main requesting new jsapi ticket")
		return "", 1
	}

	jsapiTicketMutex.Lock()
	defer jsapiTicketMutex.Unlock()
	jsapiTicket = resp.Ticket
	log.Printf("[wxtoken.cacheJSApiTicket] caching:main get response %v\n", resp)
	return resp.Ticket, resp.ExpiresIn
}

func main() {
	c := getConfig()
	appid = c.AppID
	appsecret := c.AppSecret
	tokenCH = make(chan string)
	go func() {
		log.Println("[wxtoken] caching:main start caching token")
		for {
			accessToken, t := cacheAccessToken(appid, appsecret)
			tokenCH <- accessToken
			time.Sleep(time.Second * time.Duration(t))
		}
	}()

	go func() {
		log.Println("[wxtoken] caching:main start caching jsapi ticket")
		for {
			_, t := cacheJSApiTicket(<-tokenCH)
			time.Sleep(time.Second * time.Duration(t))
		}
	}()

	e := echo.New()
	e.Use(middleware.Recover())
	e.GET("/service/resources/AccessToken", accessTokenView)
	e.GET("/service/resources/JsApiTicket", jsapiTicketView)
	e.GET("/service/resources/signature", jssdkConfigView)
	// e.Run(standard.New(c.Addr))
	e.Logger.Fatal(e.Start(c.Addr))
}
