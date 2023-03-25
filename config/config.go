package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type BotConfig struct {
	APIKey        string
	LogLevel      string
	LogOutputFile string
	LogServer     string
	ServiceName   string
	APIHolidayKey string
}

type Holiday struct {
	Name    string `json:"name"`
	Date    string `json:"date"`
	Country string `json:"country"`
}

func NewBotInfastructureConfig() *BotConfig {

	// Create a new instance of the BotConfig struct
	config := &BotConfig{}

	// Load environment variables from .env file
	loadFile := godotenv.Load(".env")
	if loadFile != nil {
		log.Fatalf("Error loading .env file: %s", loadFile.Error())
	}

	// Populate the configuration settings
	config.APIKey = os.Getenv("API_KEY")
	config.LogLevel = os.Getenv("LOG_LEVEL")
	config.LogOutputFile = os.Getenv("LOG_OUTPUT_FILE")
	config.LogServer = os.Getenv("LOG_SERVER")
	config.ServiceName = os.Getenv("SERVICE_NAME")
	config.APIHolidayKey = os.Getenv("API_KEY_HOLIDAYS")

	return config
}
