package router

import (
	"fmt"
	"io"
	"os"

	"github.com/Beretta350/authentication/internal/app/controller"
	"github.com/Beretta350/authentication/internal/app/middleware"
	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	app := gin.New()

	f, _ := os.Create("log/app.log")
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(f)
	app.SetTrustedProxies(nil)

	//middlewares
	app.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - - [%s] \"%s %s %s %d %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format("2006/01/02 15:04:05"),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
		)
	}))

	app.Use(gin.Recovery())
	app.Use(middleware.GlobalErrorHandler())

	return app
}

func SetupUserRoutes(engine *gin.Engine, controller controller.UserController) *gin.Engine {
	engine.POST("/save", controller.Save)
	engine.POST("/login", controller.Login)
	return engine
}
