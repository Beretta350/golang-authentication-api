package database

import (
	"context"
	"fmt"

	"github.com/Beretta350/authentication/config"
)

type DatabaseConfig interface {
	ConnectDB(ctx context.Context) any
	GetURI() string
}

func getDbUri() string {
	var uri string
	dbConfig := config.GetConfig().Database

	switch dbConfig.GetProtocol() {
	case "mongodb+srv":
		uri = fmt.Sprintf(
			"mongodb+srv://%s:%s@%s.mongodb.net/?retryWrites=true&w=majority&appName=%s",
			dbConfig.GetUsername(),
			dbConfig.GetPassword(),
			dbConfig.GetClusterAddress(),
			dbConfig.GetAppName(),
		)
	case "mongodb":
		uri = fmt.Sprintf(
			"mongodb://%s:%s@%s:%d/",
			dbConfig.GetUsername(),
			dbConfig.GetPassword(),
			dbConfig.GetHost(),
			dbConfig.GetPort(),
		)
	default:
		panic("Unsupported database protocol")
	}

	fmt.Printf("URI: %v\n", uri)
	return uri
}
