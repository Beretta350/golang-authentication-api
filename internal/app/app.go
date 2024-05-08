package app

import (
	"context"
	"fmt"
	"log"

	"github.com/Beretta350/authentication/config"
	"github.com/Beretta350/authentication/internal/app/common/router"
	userController "github.com/Beretta350/authentication/internal/app/user/controller"
	userRepo "github.com/Beretta350/authentication/internal/app/user/repository"
	userService "github.com/Beretta350/authentication/internal/app/user/service"
	"github.com/Beretta350/authentication/internal/pkg/database"
	"github.com/Beretta350/authentication/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func Run(env string) {
	if env == "" {
		env = "local"
	}
	config.Setup(env)

	cfg := config.GetConfig()
	ctx := context.Background()

	jwtWrap := jwt.NewJWTWrapper(cfg.JWTSecret)

	//mongodb
	mongodb := database.ConnectDB(ctx, cfg.Database)

	//repositories
	userRepo := userRepo.NewUserRepository(mongodb)

	//services
	userService := userService.NewUserService(userRepo)

	//controllers
	userController := userController.NewUserController(userService, jwtWrap)

	web := router.Setup()
	web = router.SetupUserRoutes(web, userController)

	gin.SetMode(cfg.Server.Mode)
	log.Printf("Server running on port %v in %v mode\n", cfg.Server.Port, cfg.Server.Mode)
	_ = web.Run(":" + fmt.Sprint(cfg.Server.Port))
}
