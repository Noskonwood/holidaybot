package main

import (
	_ "go.uber.org/zap"

	"holidaybot/bot"
	"holidaybot/config"
	"holidaybot/container"
	"holidaybot/logger"
)

func main() {
	// Load configuration
	botConfig := config.NewBotInfastructureConfig()

	// Initialize logger
	log, err := logger.NewBotInfrastructureLogger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logger.Close(log)

	// Initialize bot infrastructure container
	botContainer := container.NewBotInfrastructureContainer(botConfig, log)

	// Start the bot
	bot.Bot(botContainer)
}
