package main

import (
	"censys/cmd/server"
	"log"
)

// @title Censys Asset Management API
// @version 1.0
// @description API for managing IT assets, ports, and tags with automated risk assessment

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

func main() {
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
