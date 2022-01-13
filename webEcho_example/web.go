package webEcho

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/cqlsy/leolib/leoconfig"
	"github.com/cqlsy/leolib/leowebecho"
	"net/http"
)

// leogo
func StartWeb() {
	leowebecho.NewEcho()
	leowebecho.SetRetType(leowebecho.RET_JSON)
	leowebecho.Logger()
	leowebecho.Recover()
	//	leoEcho.Debug(true)
	corsConfig := middleware.CORSConfig{
		Skipper:          middleware.DefaultSkipper,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
		AllowCredentials: true,
	}
	leowebecho.CORS(corsConfig)
	// 初始化静态资源访问目录
	leowebecho.Echo.Static("/static", "view/static")
	leowebecho.Echo.Static("/", "data")
	leowebecho.Echo.HTTPErrorHandler = customHTTPErrorHandler
	initRouter()
	_ = leowebecho.Echo.Start(fmt.Sprintf("%v:%v", leoconfig.Conf.Web.Ip, leoconfig.Conf.Web.Port))
}

// CustomHTTPErrorHandler
// 自定义错误handler处理
func customHTTPErrorHandler(err error, c echo.Context) {
	var (
		code = http.StatusInternalServerError
		msg  interface{}
	)
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		msg = he.Message
	} else if c.Echo().Debug {
		msg = err.Error()
	} else {
		msg = http.StatusText(code)
	}
	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{
			"code": -1,
			"data": nil,
			"msg":  msg,
		}
	}
	if !c.Response().Committed {
		if c.Request().Method == "HEAD" { // Issue #608
			if err := c.NoContent(code); err != nil {
				c.Echo().Logger.Error(err)
			}
		} else {
			if err := c.JSON(code, msg); err != nil {
				c.Echo().Logger.Error(err)
			}
		}
	}
}

// 初始化路由
func initRouter() {
	leowebecho.Echo.GET("/test", nil)
}
