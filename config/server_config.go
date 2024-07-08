package config

import (
	"log"
	"os"
	"strconv"
)

type serverConfig struct {
	port int
	mode string
}

// TODO: Need to refactor to accept multiple servers configs
func setupServer() *serverConfig {
	serverMode := os.Getenv("SERVER_MODE")
	serverPort, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Fatal("invalid server port: %w", err)
	}

	return &serverConfig{
		port: serverPort,
		mode: serverMode,
	}
}

func (s *serverConfig) GetPort() int {
	return s.port
}

func (s *serverConfig) GetMode() string {
	return s.mode
}
