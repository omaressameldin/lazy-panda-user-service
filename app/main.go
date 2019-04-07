package main

import (
	"fmt"
	"os"

	"github.com/omaressameldin/lazy-panda-user-service/cmd"
)

func main() {
	defer cmd.CloseServer()

	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
