package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/joho/godotenv"
)

var config *Configuration
var configOnce sync.Once

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type Config interface {
	SetupConfig() *Configuration
	GetConfig() *Configuration
}

type Configuration struct {
	Database  *databaseConfig
	Server    *serverConfig
	JWTSecret string
}

// singleton
func SetupConfig(configFile string) *Configuration {
	configOnce.Do(func() {
		configPath := filepath.Join(basepath, configFile+".env")
		err := godotenv.Load(configPath)
		if err != nil && configFile == "local" {
			log.Fatalln("error loading environment file: %w", err)
		}

		config = &Configuration{
			Database:  setupDatabase(),
			Server:    setupServer(),
			JWTSecret: os.Getenv("JWT_SECRET"),
		}
	})

	return config
}

func GetConfig() *Configuration {
	return config
}
