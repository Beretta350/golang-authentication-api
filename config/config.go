package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/joho/godotenv"
)

var config *Configuration

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type ServerConfig struct {
	Port int
	Mode string
}

type Config interface {
	Setup(configFile string)
	GetConfig() *Configuration
}

type Configuration struct {
	Database   DatabaseConfig
	Server     ServerConfig
	JWTSecret  string
	CSRFSecret string
}

func Setup(configFile string) {
	configPath := filepath.Join(basepath, configFile+".env")
	err := godotenv.Load(configPath)
	if err != nil && configFile == "local" {
		log.Fatalln("error loading environment file: %w", err)
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalln("invalid database port: %w", err)
	}

	dbConfig := DatabaseConfig{
		Driver:   os.Getenv("DB_DRIVER"),
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_DATABASE"),
	}

	// Server configuration
	serverMode := os.Getenv("SERVER_MODE")
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("invalid server port: %w", err)
	}

	serverConfig := ServerConfig{
		Port: serverPort,
		Mode: serverMode,
	}

	config = &Configuration{
		Database:   dbConfig,
		Server:     serverConfig,
		JWTSecret:  os.Getenv("JWT_SECRET"),
		CSRFSecret: os.Getenv("CSRF_SECRET"),
	}
}

func GetConfig() *Configuration {
	return config
}
