package config

import "fmt"

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Username string
	Password string
	Database string
	Sslmode  bool
}

func (db *DatabaseConfig) GetURI() string {

	var uri string

	switch db.Driver {
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
