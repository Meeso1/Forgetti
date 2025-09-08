package main

import (
	"ForgettiServer/config"
	"ForgettiServer/routes"
	"ForgettiServer/services"
	"fmt"
	"log"
	"path/filepath"

	"forgetti-common/logging"

	"github.com/gin-gonic/gin"
)

const configFile = "config.json"

func main() {
	cfg, err := config.LoadFromFile(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	setupLogging(cfg)
	
	logger := logging.MakeLogger("main")
	logger.Verbose("Logging configured successfully")
	logger.Info("Starting Forgetti Server...")

	gin.SetMode(cfg.Server.Mode)
	logger.Verbose("Gin mode set to: %s", cfg.Server.Mode)

	r := gin.Default()
	logger.Verbose("Gin router initialized")

	logger.Verbose("Creating service container...")
	serviceContainer, err := services.CreateServiceContainer(cfg)
	if err != nil {
		logger.Error("Failed to create service container: %v", err)
		log.Fatalf("Failed to create service container: %v", err)
	}
	logger.Verbose("Service container created successfully")

	logger.Verbose("Setting up routes...")
	routes.AddEncRoutes(r, serviceContainer)
	logger.Verbose("Routes configured successfully")

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	logger.Info("Starting server on %s in %s mode", addr, cfg.Server.Mode)
	if err := r.Run(addr); err != nil {
		logger.Error("Failed to start server: %v", err)
	}
}

func setupLogging(cfg *config.Config) {
	logging.SetGlobalConfig(logging.Config{
		LogLevel: logging.LogLevelInfo,
		LogFile:  "",
	})
	logger := logging.MakeLogger("main")

	logLevel := logging.LogLevelInfo
	switch cfg.Logging.Level {
	case "debug":
		logLevel = logging.LogLevelVerbose
	case "info":
		logLevel = logging.LogLevelInfo
	case "warn", "error":
		logLevel = logging.LogLevelError
	}

	logFile := ""
	if cfg.Logging.LogFile != "" {
		logFile = filepath.Join(cfg.Logging.LogDirectory, cfg.Logging.LogFile)
		logger.Info("Configuring file logging to: %s", logFile)
	}

	err := logging.SetGlobalConfig(logging.Config{
		LogLevel: logLevel,
		LogFile:  logFile,
	})
	if err != nil {
		logger.Error("Failed to configure logging: %v", err)
	}
}
