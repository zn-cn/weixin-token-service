package daemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"util"
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
)
var ticketLog = util.GetLogger("/app/log/daemon/ticket.txt", "[DEBUG]")

// Update 更新任务
func (h *Handler) Update() {
	// access_token 和 jsapi_ticket 应该同时更新
	h.UpdateToekn()
	h.UpdateTicket()
}

// UpdateTicket 更新ticket
func (h *Handler) UpdateTicket() {
	accessToken, err := (*(h.RedisConn[0])).Do("GET", AccessToken)
	if err != nil {
		ticketLog.Println("access_token 缺失")
		return
	}

	url := fmt.Sprintf(model.WeixinURL[JSApiTicket], accessToken)
	resp, err := http.Get(url)
	if err != nil {
		tokenLog.Println("获取jsapi_ticket失败")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tokenLog.Println("json 解析失败")
		return
	}

	ticket, ok := result["ticket"]
	if !ok {
		tokenLog.Println(fmt.Sprintf("errcode:%s, errmsg:%s", result["errcode"], result["errmsg"]))
		return
	}

	for i := 0; i < 3; i++ {
		(*(h.RedisConn[i])).Do("SETEX", JSApiTicket, int(result["expires_in"].(float64)), ticket.(string))
	}
}
