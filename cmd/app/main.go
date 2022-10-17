package main

import (
	"log"

	"alukart32.com/bank/config"
	"alukart32.com/bank/internal/app"
)

// The entry point of the application.
func main() {
	cfg, err := config.New(config.Default)
	if err != nil {
		log.Fatal("read config error: ", err)
	}

	app.Run(*cfg)
}
