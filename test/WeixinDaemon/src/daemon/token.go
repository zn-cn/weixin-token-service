package daemon

import (
	"conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"model"
	"net/http"
	"util"
)

var tokenLog = util.GetLogger("/app/log/daemon/token.txt", "[DEBUG]")

// UpdateToekn 更新access_token
func (h *Handler) UpdateToekn() {
	url := fmt.Sprintf(model.WeixinURL["access_token"], conf.Conf.Weixin.AppID, conf.Conf.Weixin.AppSecret)
	resp, err := http.Get(url)
	if err != nil {
		tokenLog.Println("获取access_token失败")
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		tokenLog.Println("json 解析失败")
		return
	}
	accessToken, ok := result["access_token"]
	if !ok {
		tokenLog.Println(fmt.Sprintf("errcode:%s, errmsg:%s", result["errcode"], result["errmsg"]))
		return
	}
	for i := 0; i < 3; i++ {
		(*(h.RedisConn[i])).Do("SETEX", "access_token", int(result["expires_in"].(float64)), accessToken.(string))
	}
}
