package util

import (
	"fmt"
	"os"
)

var commonLog = GetLogger("/app/log/util/common.txt", "[DEBUG]")

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
