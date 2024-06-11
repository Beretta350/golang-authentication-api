package config

import "fmt"

type DatabaseConfig struct {
	Driver         string
	Host           string
	Port           int
	Username       string
	Password       string
	Database       string
	Sslmode        bool
	ClusterAddress string
	AppName        string
}

func (db *DatabaseConfig) GetURI() string {

	var uri string

	switch db.Driver {
	case "mongodb+srv":
		uri = fmt.Sprintf("mongodb+srv://%s:%s@%s.mongodb.net/?retryWrites=true&w=majority&appName=%s",
			db.Username, db.Password, db.ClusterAddress, db.AppName)
	case "mongodb":
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/",
			db.Username, db.Password, db.Host, db.Port)
	default:
		panic("Unsupported database driver")
	}

	fmt.Printf("URI: %v\n", uri)

	return uri
}
