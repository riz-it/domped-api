package main

import (
	"fmt"
	"log"
	"strconv"

	"riz.it/domped/internal/config"
	"riz.it/domped/internal/injector"
)

func main() {
	app := injector.InitializedApp()
	cnf := config.Get()

	port, _ := strconv.Atoi(cnf.Server.Port)
	err := app.Fiber.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
