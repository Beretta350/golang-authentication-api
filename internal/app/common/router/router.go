package router

import (
	"fmt"
	"io"
	"os"

	"github.com/Beretta350/authentication/config"
	router_constants "github.com/Beretta350/authentication/internal/app/common/constants/router"
	"github.com/Beretta350/authentication/internal/app/common/middleware"
	userController "github.com/Beretta350/authentication/internal/app/user/controller"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const localhostPath string = "http://localhost:8080"

func Setup(cfg *config.Configuration) *gin.Engine {
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

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{localhostPath}

	jwtWrap := jwt.NewJWTWrapper(
		cfg.JWTSecret,
		[]string{router_constants.LoginRoute, router_constants.SaveRoute},
	)

	app.Use(gin.Recovery())
	app.Use(cors.New(corsConfig))
	app.Use(middleware.JWTHandler(jwtWrap))
	app.Use(middleware.GlobalErrorHandler())

	gin.SetMode(cfg.Server.Mode)
	return app
}

func SetupUserRoutes(engine *gin.Engine, controller userController.UserController) *gin.Engine {
	engine.GET(router_constants.RefreshTokenRoute, controller.RefreshToken)
	engine.GET(router_constants.GetUserByIDRoute, controller.GetUserByID)
	engine.POST(router_constants.LoginRoute, controller.Login)
	engine.POST(router_constants.SaveRoute, controller.Save)
	engine.PUT(router_constants.UpdateRoute, controller.Update)
	engine.DELETE(router_constants.DeleteRoute, controller.Delete)
	return engine
}
