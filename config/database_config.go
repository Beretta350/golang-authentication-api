package config

import (
	"log"
	"os"
	"strconv"
)

// TODO: Need to refactor to accept multiple db configs
func setupDatabase() *databaseConfig {
	var dbPort int = 0
	var err error = nil
	if os.Getenv("DB_PORT") != "" {
		dbPort, err = strconv.Atoi(os.Getenv("DB_PORT"))
		if err != nil {
			log.Fatalln("invalid database port: %w", err)
		}
	}

	return &databaseConfig{
		protocol:       os.Getenv("DB_PROTOCOL"),
		host:           os.Getenv("DB_HOST"),
		port:           dbPort,
		username:       os.Getenv("DB_USERNAME"),
		password:       os.Getenv("DB_PASSWORD"),
		database:       os.Getenv("DB_DATABASE"),
		clusterAddress: os.Getenv("DB_CLUSTER_ADDRS"),
		appName:        os.Getenv("DB_APPNAME"),
	}
}

type databaseConfig struct {
	protocol       string
	host           string
	port           int
	username       string
	password       string
	database       string
	sslmode        bool
	clusterAddress string
	appName        string
}

func (db *databaseConfig) GetProtocol() string {
	return db.protocol
}

func (db *databaseConfig) GetHost() string {
	return db.host
}

func (db *databaseConfig) GetPort() int {
	return db.port
}

func (db *databaseConfig) GetUsername() string {
	return db.username
}

func (db *databaseConfig) GetPassword() string {
	return db.password
}

func (db *databaseConfig) GetDatabase() string {
	return db.database
}

func (db *databaseConfig) GetSslMode() bool {
	return db.sslmode
}

func (db *databaseConfig) GetClusterAddress() string {
	return db.clusterAddress
}

func (db *databaseConfig) GetAppName() string {
	return db.appName
}
