package routes

import (
	"conf"
	"fmt"
	"handler"
	"util"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var routesLogger = util.GetLogger("/app/restful/routes/routes.txt", "[DEBUG]")

// Init 初始化
func Init() {

	e := echo.New()
	e.Use(middleware.Logger())

	routes(e)

	// Start server
	addr := fmt.Sprintf("%s:%s", conf.Conf.App.Host, conf.Conf.App.Port)
	e.Logger.Fatal(e.Start(addr))

}

func routes(e *echo.Echo) {
	e.GET("/health", handler.HelloWorld)
	g := e.Group("/service/resources")
	g.GET("/AccessToken", handler.GetToken)
	g.GET("/JsApiTicket", handler.GetTicket)
	g.GET("/signature", handler.GetSignature)
}
