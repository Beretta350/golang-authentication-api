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
)

func Run(env string) {
	if env == "" {
		env = "local"
	}

	ctx := context.Background()
	cfg := config.SetupConfig(env)

	//mongodb
	dynamodbConn := database.NewDynamoDBSession().ConnectDB(ctx)

	//repositories
	userRepo := userRepo.NewDynamoDBUserRepository(dynamodbConn)

	//services
	userService := userService.NewUserService(userRepo)

	//controllers
	userController := userController.NewUserController(userService)

	//routes
	web := router.Setup(cfg)
	web = router.SetupUserRoutes(web, userController)

	//run
	log.Printf("Server running on port %v in %v mode\n", cfg.Server.GetPort(), cfg.Server.GetMode())
	_ = web.Run(":" + fmt.Sprint(cfg.Server.GetPort()))
}
