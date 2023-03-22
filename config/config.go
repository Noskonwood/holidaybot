package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type BotConfig struct {
	APIKey        string
	LogLevel      string
	LogOutputFile string
	LogServer     string
	ServiceName   string
}

func NewBotInfastructureConfig() *BotConfig {

	// Create a new instance of the BotConfig struct
	config := &BotConfig{}

	// Load environment variables from .env file
	loadFile := godotenv.Load("config/.env")
	if loadFile != nil {
		log.Fatalf("Error loading .env file: %s", loadFile.Error())
	}

	// Populate the configuration settings
	config.APIKey = os.Getenv("API_KEY")
	config.LogLevel = os.Getenv("LOG_LEVEL")
	config.LogOutputFile = os.Getenv("LOG_OUTPUT_FILE")
	config.LogServer = os.Getenv("LOG_SERVER")
	config.ServiceName = os.Getenv("SERVICE_NAME")

	return config
}
