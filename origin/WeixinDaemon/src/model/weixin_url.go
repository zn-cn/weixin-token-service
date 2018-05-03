package model

// WeixinURL 微信URL
var WeixinURL = map[string]string{
	"access_token": "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
	"jsapi_ticket": "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi",
}
