package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/omaressameldin/lazy-panda-user-service/server"
)

// Config is configuration for Server
type Config struct {
	Port string
}

// RunServer get the flags and starts the server
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.Port, "port", "", "port to bind")
	flag.Parse()

	if len(cfg.Port) == 0 {
		return fmt.Errorf("invalid TCP port for server: '%s'", cfg.Port)
	}

	// v1API := v1.NewToDoServiceServer(cfg.DatabasePath)

	return server.RunServer(ctx, cfg.Port)
}

// CloseServer closes all connections such as database connection
func CloseServer() error {
	return nil
	// return v1.CloseConnection()
}
