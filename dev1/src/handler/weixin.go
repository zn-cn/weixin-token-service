package handler

import (
	"cache"
	"conf"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"time"
	"util"

	"github.com/labstack/echo"
)

// GetToken 获取access_token
func GetToken(c echo.Context) (err error) {
	// 读锁
	cache.TokenMutex.RLock()
	defer cache.TokenMutex.RUnlock()

	return c.JSON(http.StatusOK, map[string]string{"access_token": cache.AccessToken})
}

// GetTicket 获取ticket
func GetTicket(c echo.Context) (err error) {

	cache.JSApiTicketMutex.RLock()
	defer cache.JSApiTicketMutex.RUnlock()
	return c.JSON(http.StatusOK, map[string]string{"jsapi_ticket": cache.JSApiTicket})
}

// GetSignature 获取signature
func GetSignature(c echo.Context) (err error) {

	cache.JSApiTicketMutex.RLock()
	defer cache.JSApiTicketMutex.RUnlock()

	url := c.QueryParam("url")
	result := getJSApiSignature(cache.JSApiTicket, url)
	result["appId"] = conf.Conf.Weixin.AppID

	return c.JSON(http.StatusOK, result)
}

// 得到用于JS_SDK的signature
func getJSApiSignature(ticket, url string) map[string]string {
	timestamp := fmt.Sprint(time.Now().Unix())
	noncestr := util.GetRandStr(16)
	tmpStr := `jsapi_ticket=` + ticket +
		`&noncestr=` + noncestr +
		`&timestamp=` + timestamp +
		`&url=` + url
	s := sha1.New()
	io.WriteString(s, tmpStr)
	signature := fmt.Sprintf("%x", s.Sum(nil))
	return map[string]string{"signature": signature, "nonce_str": noncestr, "timestamp": timestamp}
}
