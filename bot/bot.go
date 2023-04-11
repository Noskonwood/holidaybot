package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"

	"holidaybot/config"
	"holidaybot/container"
)

func Bot(container container.BotInfastructureContainer) {
	cfg := container.GetConfig()

	bot, err := tgbotapi.NewBotAPI(cfg.APIKey)
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

		switch update.Message.Text {
		case "/start":
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a country by entering its 2-letter ISO country code (e.g. JP for Japan):")
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			if _, err := bot.Send(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
				log.Fatalf("Error to send the message: %v", err)
			}

		default:
			apiKey := cfg.APIHolidayKey
			holidayAPIURL := fmt.Sprintf("https://app.abstractapi.com/api/holidays?api_key=%s&country=%s&year=%d&month=%d&day=%d", apiKey, update.Message.Text, time.Now().Year(), time.Now().Month(), time.Now().Day())

			resp, err := http.Get(holidayAPIURL)
			if err != nil {
				log.Fatalf("Error sending request to API: %v", err)
			}
			defer resp.Body.Close()

			var holidays []config.Holiday
			if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
				log.Fatalf("Error decoding response body: %v", err)
			}

			var holidayMessage string
			if len(holidays) > 0 {
				holidayMessage = fmt.Sprintf("Today is %s in %s", holidays[0].Name, update.Message.Text)
			} else {
				holidayMessage = fmt.Sprintf("There are no holidays today in %s", update.Message.Text)
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, holidayMessage)

			if _, err := bot.Send(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
				log.Fatalf("Error to send the message: %v", err)
			}
		}
	}
}
