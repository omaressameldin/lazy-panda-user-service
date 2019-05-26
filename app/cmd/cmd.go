package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	v1 "github.com/omaressameldin/lazy-panda-user-service/app/pkg/service/v1"
	"github.com/omaressameldin/lazy-panda-user-service/app/server"
	"github.com/omaressameldin/lazy-panda-utils/app/pkg/database"
	"github.com/omaressameldin/lazy-panda-utils/app/pkg/firebase"
)

// Config is configuration for Server
type Config struct {
	Port           string
	FirebaseConfig string
	Collection     string
	Bucket         string
}

var v1API *v1.UserServiceServer

// RunServer get the flags and starts the server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Port, "port", "", "port to bind")
	flag.StringVar(&cfg.FirebaseConfig, "firebaseConfig", "", "firebase json config file")
	flag.StringVar(&cfg.Collection, "collection", "", "firebase collection")
	flag.StringVar(&cfg.Bucket, "bucket", "", "firebase storage bucket")
	flag.Parse()

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	if len(cfg.Collection) == 0 {
		return fmt.Errorf("invalid Collection for firebase database: '%s'", cfg.Collection)
	}

	if len(cfg.Bucket) == 0 {
		return fmt.Errorf("invalid Collection for firebase database: '%s'", cfg.Collection)
	}

	_, err := os.Stat(cfg.FirebaseConfig)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: '%s'", cfg.FirebaseConfig)
	}
	connector := initConnector(cfg.FirebaseConfig, cfg.Collection, cfg.Bucket)
	v1API = v1.NewUserServiceServer(connector)

	return server.RunServer(ctx, v1API, cfg.Port)
}

// initConnector initializes database connector
func initConnector(firebaseConfig, collection, bucket string) database.Connector {
	connector, err := firebase.StartConnection(firebaseConfig, collection, bucket)
	if err != nil {
		panic(err)
	}

	return connector
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return v1API.CloseConnection()
}
