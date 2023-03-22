package main

import (
	"git.foxminded.ua/foxstudent104181/holidaybot/bot"
	"git.foxminded.ua/foxstudent104181/holidaybot/config"
	"git.foxminded.ua/foxstudent104181/holidaybot/container"
	"git.foxminded.ua/foxstudent104181/holidaybot/logger"
	_ "go.uber.org/zap"
)

func main() {
	// Initialize logger
	log, err := logger.NewBotInfrastructureLogger("DEBUG")
	if err != nil {
		panic(err)
	}
	defer logger.Close(log)

	// Load configuration
	botConfig := config.NewBotInfastructureConfig()

	// Initialize bot infrastructure container
	container := container.NewBotInfrastructureContainer(botConfig, log)

	// Start the bot
	bot.Bot(container)
}
