package main

import (
	"fmt"
	"log"

	"alukart32.com/bank/config"
	"alukart32.com/bank/internal/app"
)

// The entry point of the application.
func main() {
	cfg, err := config.New(config.Default)
	if err != nil {
		log.Fatal(fmt.Errorf("read config error: %w", err))
	}

	app.Run(cfg)
}
