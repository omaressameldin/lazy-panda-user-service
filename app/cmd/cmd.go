package cmd

import (
	"context"
	"flag"
	"fmt"
	"os"

	v1 "github.com/omaressameldin/lazy-panda-user-service/pkg/service/v1"
	"github.com/omaressameldin/lazy-panda-user-service/server"
)

// Config is configuration for Server
type Config struct {
	Port           string
	FirebaseConfig string
	Collection     string
}

// RunServer get the flags and starts the server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Port, "port", "", "port to bind")
	flag.StringVar(&cfg.FirebaseConfig, "firebaseConfig", "", "firebase json config file")
	flag.StringVar(&cfg.Collection, "collection", "", "firebase collection")

	flag.Parse()

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	if len(cfg.Collection) == 0 {
		return fmt.Errorf("invalid Collection for firebase database: '%s'", cfg.Collection)
	}

	_, err := os.Stat(cfg.FirebaseConfig)
	if os.IsNotExist(err) {
		return fmt.Errorf("File does not exist: '%s'", cfg.FirebaseConfig)
	}

	v1API := v1.NewUserServiceServer(cfg.FirebaseConfig, cfg.Collection)

	return server.RunServer(ctx, v1API, cfg.Port)
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return v1.CloseConnection()
}
