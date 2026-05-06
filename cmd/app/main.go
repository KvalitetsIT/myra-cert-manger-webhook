package main

import (
	"log/slog"

	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/configs"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/logging"
	"github.com/KvalitetsIT/cert-manager-webhook-core/pkg/service"
	"github.com/KvalitetsIT/cert-manager-webhook-myra/internal/client"
	internal "github.com/KvalitetsIT/cert-manager-webhook-myra/internal/client/adaptors"
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
}

func main() {
	myra, err := client.NewMyraClient(CONFIG.Myra, logger)
	if err != nil {
		panic(err)
	}

	adaptor := internal.NewMyraClientAdaptor(myra)

	factory := service.NewServiceFactory(CONFIG, logger)
	if service, err := factory.CreateDefault(adaptor); err != nil {
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
