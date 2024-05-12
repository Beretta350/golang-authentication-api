package router

import (
	"fmt"
	"io"
	"os"

	"github.com/Beretta350/authentication/config"
	"github.com/Beretta350/authentication/internal/app/common/middleware"
	userController "github.com/Beretta350/authentication/internal/app/user/controller"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
	corsConfig.AllowOrigins = []string{"http://localhost:8080"}

	jwtWrap := jwt.NewJWTWrapper(cfg.JWTSecret, []string{"/login", "/save", "/refresh"})

	app.Use(gin.Recovery())
	app.Use(cors.New(corsConfig))
	app.Use(middleware.JWTHandler(jwtWrap))
	app.Use(middleware.GlobalErrorHandler())

	gin.SetMode(cfg.Server.Mode)
	return app
}

func SetupUserRoutes(engine *gin.Engine, controller userController.UserController) *gin.Engine {
	engine.GET("/refresh", controller.RefreshToken)
	engine.GET("/user", controller.GetUserByID)
	engine.POST("/login", controller.Login)
	engine.POST("/save", controller.Save)
	engine.PUT("/update", controller.Update)
	engine.DELETE("/delete", controller.Delete)
	return engine
}
