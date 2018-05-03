package handler

import (
	"conf"
	"crypto/sha1"
	"fmt"
	"io"
	"net/http"
	"time"
	"util"

	"github.com/garyburd/redigo/redis"

	"github.com/labstack/echo"
)

var (
	// AccessToken : redis 中 微信 access_token 的 key
	AccessToken = "access_token"
	// JSApiTicket : redis 中 微信 jsapi_ticket 的 key
	JSApiTicket = "jsapi_ticket"
	// JSApiNoncestr : redis中 微信JS-sdk的 noncestr 的 key (一个随机字符串)
	JSApiNoncestr = "jsapi_noncestr"
	// JSApiTimestamp : redis中 微信JS-sdk的 timestamp key  (一个随机字符串)
	JSApiTimestamp = "jsapi_timestamp"
	// JSApiSignature : redis中 微信JS-sdk的 signature key  (一个随机字符串)
	JSApiSignature = "jsapi_signature"
	weixinLog      = util.GetLogger("/app/log/handler/weixin.txt", "[DEBUG]")
)

// GetToken 获取access_token
func (h *Handler) GetToken(c echo.Context) (err error) {
	redisConn := h.RedisConn.Get()
	defer redisConn.Close()
	accessToken, err := redis.String(redisConn.Do("GET", AccessToken))
	if err != nil {
		weixinLog.Println("access_token 获取失败")
		return c.JSON(http.StatusInternalServerError, map[string]string{"errmsg": "access_token 获取失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"access_token": accessToken})
}

// GetTicket 获取ticket
func (h *Handler) GetTicket(c echo.Context) (err error) {
	redisConn := h.RedisConn.Get()
	defer redisConn.Close()
	ticket, err := redis.String(redisConn.Do("GET", JSApiTicket))
	if err != nil {
		weixinLog.Println("jsapi_ticket 获取失败")
		return c.JSON(http.StatusInternalServerError, map[string]string{"errmsg": "jsapi_ticket 获取失败"})
	}
	return c.JSON(http.StatusOK, map[string]string{"jsapi_ticket": ticket})
}

// GetSignature 获取signature
func (h *Handler) GetSignature(c echo.Context) (err error) {
	redisConn := h.RedisConn.Get()
	defer redisConn.Close()
	ticket, err := redis.String(redisConn.Do("GET", JSApiTicket))
	if err != nil {
		weixinLog.Println("jsapi_ticket 获取失败")
		return c.JSON(http.StatusInternalServerError, map[string]string{"errmsg": "jsapi_ticket 获取失败"})
	}

	url := c.QueryParam("url")
	result := getJSApiSignature(ticket, url)
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
