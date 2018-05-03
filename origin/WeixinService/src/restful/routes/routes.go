package routes

import (
	"conf"
	"fmt"
	"handler"
	"restful"
	"util"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var routesLogger = util.GetLogger("/app/restful/routes/routes.txt", "[DEBUG]")

// Init 初始化
func Init() {

	e := echo.New()

	// InitConfig
	if err := conf.InitConfig(""); err != nil {
		routesLogger.Fatalln(err)
	}

	e.Use(middleware.Logger())

	// Database connection
	h, err := restful.DbInit()
	defer (*h.RedisConn).Close()
	if err != nil {
		routesLogger.Fatalln(err)
	}

	routes(e, h)

	// Start server
	addr := fmt.Sprintf("%s:%s", conf.Conf.App.Host, conf.Conf.App.Port)
	e.Logger.Fatal(e.Start(addr))

}

func routes(e *echo.Echo, h *handler.Handler) {
	e.GET("/health", handler.HelloWorld)
	g := e.Group("/service/resources")
	g.GET("/AccessToken", h.GetToken)
	g.GET("/JsApiTicket", h.GetTicket)
	g.GET("/signature", h.GetSignature)
}
