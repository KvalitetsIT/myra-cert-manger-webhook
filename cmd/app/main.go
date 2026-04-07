package main

import (
	"log/slog"

	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/configs"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/logging"
	"github.com/KvalitetsIT/myra-cert-manager-webhook/internal/service"
	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

var CONFIG configs.Configuration
var logger *slog.Logger

func init() {

	logger = logging.NewJSONLogger()
	slog.SetDefault(logger)

	loadEnv()
	parseEnv()

	logger.Info("Loaded configuration", slog.Any("Configuration", CONFIG))
}

func main() {
	factory := service.NewServiceFactory(CONFIG, logger)
	if service, err := factory.CreateDefault(); err != nil {
		logger.Error("Failed to create default service", slog.Any("error", err))
	} else {
		service.Start()
	}
}

// Parse env into struct
func parseEnv() {
	if err := env.Parse(&CONFIG); err != nil {
		logger.Error("Could not parse environment variables", slog.Any("error", err))
	}
}

// Load .env file
func loadEnv() {
	if err := godotenv.Load(); err != nil {
		logger.Error(".env file not found, using system environment variables")
	}
}
