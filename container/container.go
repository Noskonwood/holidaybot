package container

import (
	"git.foxminded.ua/foxstudent104181/holidaybot/config"
	"go.uber.org/zap"
)

// BotInfrastructureContainer represents an interface for accessing the data which sharing in overall application

type BotInfastructureContainer interface {
	GetConfig() *config.BotConfig
	GetLogger() *zap.SugaredLogger
}

// container struct is for sharing data such as the setting of application and logger in overall this application
type container struct {
	config *config.BotConfig
	logger *zap.SugaredLogger
}

// NewBotInfrastructureContainer is constructor
func NewBotInfrastructureContainer(config *config.BotConfig, logger *zap.SugaredLogger) BotInfastructureContainer {
	return &container{
		config: config,
		logger: logger,
	}
}

// GetConfig returns the object of configuration
func (c *container) GetConfig() *config.BotConfig {
	return c.config
}

// GetLogger returns the object of logger
func (c *container) GetLogger() *zap.SugaredLogger {
	return c.logger
}
