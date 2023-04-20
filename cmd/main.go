package main

import (
	"example.com/holidaybot/bot"
	"example.com/holidaybot/config"
	"example.com/holidaybot/container"
	"example.com/holidaybot/logger"
	_ "go.uber.org/zap"
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
