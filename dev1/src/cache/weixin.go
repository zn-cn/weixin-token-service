package cache

import (
	"conf"
	"fmt"
	"sync"
	"util"
)

var (
	// AccessToken :  微信 access_token
	AccessToken = "access_token"
	// JSApiTicket :  微信 jsapi_ticket
	JSApiTicket = "jsapi_ticket"

	ticketLog = util.GetLogger("/app/log/web/cache/ticket.txt", "[DEBUG]")
	tokenLog  = util.GetLogger("/app/log/web/cache/token.txt", "[DEBUG]")

	// 读写锁
	TokenMutex       sync.RWMutex
	JSApiTicketMutex sync.RWMutex
	tokenCH          = make(chan string)

	// WeixinURL 微信URL
	WeixinURL = map[string]string{
		"access_token": "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		"jsapi_ticket": "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
	}
)

// UpdateToekn 更新access_token
func UpdateToekn() {
	url := fmt.Sprintf(WeixinURL["access_token"], conf.Conf.Weixin.AppID, conf.Conf.Weixin.AppSecret)

	// 写锁
	TokenMutex.Lock()
	defer TokenMutex.Unlock()

	result := util.GetJSONBody(url)
	if result == nil {
		tokenLog.Println("access_token 更新失败")
		return
	}
	accessToken, ok := result["access_token"]
	if !ok {
		tokenLog.Println(fmt.Sprintf("errcode:%s, errmsg:%s", result["errcode"], result["errmsg"]))
		return
	}
	AccessToken = accessToken.(string)
	tokenCH <- accessToken.(string)
}

// UpdateTicket 更新ticket
func UpdateTicket() {
	// 管道保证先更新token后更新ticket
	token := <-tokenCH
	url := fmt.Sprintf(WeixinURL[JSApiTicket], token)

	// 写锁
	JSApiTicketMutex.Lock()
	defer JSApiTicketMutex.Unlock()
	result := util.GetJSONBody(url)
	if result == nil {
		tokenLog.Println("获取jsapi_ticket失败")
		return
	}

	ticket, ok := result["ticket"]
	if !ok {
		tokenLog.Println(fmt.Sprintf("errcode:%s, errmsg:%s", result["errcode"], result["errmsg"]))
		return
	}
	JSApiTicket = ticket.(string)

}
