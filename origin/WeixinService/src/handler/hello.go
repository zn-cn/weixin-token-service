package handler

import (
	"net/http"

	"github.com/labstack/echo"
)

// HelloWorld 测试健康
func HelloWorld(c echo.Context) (err error) {
	return c.String(http.StatusOK, "hello world")
}
