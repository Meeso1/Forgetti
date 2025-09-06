package main

import (
	"ForgettiServer/config"
	"ForgettiServer/routes"
	"ForgettiServer/services"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

const configFile = "config.json"

// TODO: add logger (maybe extend console logger?)
func main() {
	cfg, err := config.LoadFromFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	gin.SetMode(cfg.Server.Mode)

	r := gin.Default()

	serviceContainer := services.CreateServiceContainer(cfg)
	routes.AddEncRoutes(r, serviceContainer)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Starting server on %s in %s mode", addr, cfg.Server.Mode)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
