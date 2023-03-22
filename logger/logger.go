package logger

import (
	"git.foxminded.ua/foxstudent104181/holidaybot/config"
	"go.uber.org/zap"
	"log"
)

const (
	developmentLevel = "DEBUG"
	productionLevel  = "PRODUCTION"
)

func NewBotInfrastructureLogger(string) (*zap.SugaredLogger, error) {
	var l *zap.Logger
	var err error

	switch config.NewBotInfastructureConfig().LogLevel {
	case productionLevel:
		l, err = zap.NewProduction()
		if err != nil {
			return nil, err
		}
	case developmentLevel:
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	default:
		l, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	sugar := l.Sugar()

	return sugar, nil
}

func Close(l *zap.SugaredLogger) {
	err := l.Sync()
	if err != nil {
		log.Println("Couldn't flush logging buffer")
	}
}
