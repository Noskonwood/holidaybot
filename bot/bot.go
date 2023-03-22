package bot

import (
	"git.foxminded.ua/foxstudent104181/holidaybot/config"
	"git.foxminded.ua/foxstudent104181/holidaybot/container"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

func Bot(container.BotInfastructureContainer) {
	botConfig := config.NewBotInfastructureConfig()
	bot, err := tgbotapi.NewBotAPI(botConfig.APIKey)
	if err != nil {
		log.Fatalf("Error to connect to bot: %v", err)
	}
	bot.Debug = true

	updateConfig := tgbotapi.NewUpdate(0)

	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Create a logger instance for this request
		logger := logrus.WithFields(logrus.Fields{
			"chat_id":  update.Message.Chat.ID,
			"user_id":  update.Message.From.ID,
			"username": update.Message.From.UserName,
		})

		logger.Infof("Received message: %s", update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		switch strings.ToLower(update.Message.Text) {
		case "/start", "/help":
			msg.Text = "Hi, I'm bot created by Bogdan Petrukhin\n Available commands:\n/about - provides short info about me\n/links - provides a list of my social links (GitHub, LinkedIn, etc)"
			msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
				tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("/about"),
					tgbotapi.NewKeyboardButton("/links"),
				),
			)
		case "/about":
			msg.Text = "Hi, my name is Bogdan and I'm going to be a software engineer in a future. I live in Canada and warm welcome to my telegram bot"
		case "/links":
			msg.Text = "You can find me on the following platforms:\n\nGitHub: https://github.com/Noskonwood\nLinkedIn: https://www.linkedin.com/in/bogdan-petrukhin/"
		default:
			logger.Warnf("Received unknown command: %s", update.Message.Text)
			continue
		}

		if _, err := bot.Send(msg); err != nil {
			logger.Errorf("Failed to send message: %v", err)
			log.Fatalf("Error to send the message: %v", err)
		}

		logger.Info("Message sent")
	}
}
