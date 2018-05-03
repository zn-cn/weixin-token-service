package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var (
	codes     = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/~!@#$%^&*()_="
	codeLen   = len(codes)
	commonLog = GetLogger("/app/log/util/common.txt", "[DEBUG]")
)

// GetRandStr 生成随机字符串
func GetRandStr(len int) string {
	data := make([]byte, len)
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < len; i++ {
		idx := rand.Intn(codeLen)
		data[i] = byte(codes[idx])
	}

	return string(data)
}

// TestOSENV 测试是否存在环境变量
func TestOSENV(envList []string) bool {
	for _, v := range envList {
		if os.Getenv(v) == "" {
			commonLog.Println(fmt.Sprintf("do not have %s OS_ENV", v))
			return false
		}
	}
	return true
}

// GetJSONBody 返回get请求的json数据
func GetJSONBody(url string) map[string]interface{} {

	resp, err := http.Get(url)
	if err != nil {
		commonLog.Println("获取access_token失败")
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	result := map[string]interface{}{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		commonLog.Println("json 解析失败")
		return nil
	}
	return result
}
