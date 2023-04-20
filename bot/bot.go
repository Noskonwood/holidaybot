package bot

import (
	"encoding/json"
	"example.com/holidaybot/config"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"example.com/holidaybot/container"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// Holiday represents a holiday returned by the API
type Holiday struct {
	Name string `json:"name"`
	Date string `json:"date"`
}

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
			// Show the user 2 countries to choose from with their respective flags
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a country by selecting its flag:")
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("🇺🇸 USA", "USA"),
					tgbotapi.NewInlineKeyboardButtonData("🇬🇧 UK", "UK"),
				),
			)

			if _, err := bot.Send(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
				log.Fatalf("Error to send the message: %v", err)
			}

		default:

			var country string
			if update.CallbackQuery != nil {
				switch update.CallbackQuery.Data {
				case "USA":
					country = "US"
				case "UK":
					country = "GB"
				}
			}
			if update.Message == nil || update.Message.Text != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a country by selecting its flag, don't type it :")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🇺🇸 USA", "USA"),
						tgbotapi.NewInlineKeyboardButtonData("🇬🇧 UK", "UK"),
					),
				)
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
				continue
			}
			// if user doesn't type anything
			if country == "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Choose a country by selecting its flag, don't type it :")
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(
					tgbotapi.NewInlineKeyboardRow(
						tgbotapi.NewInlineKeyboardButtonData("🇺🇸 USA", "USA"),
						tgbotapi.NewInlineKeyboardButtonData("🇬🇧 UK", "UK"),
					),
				)
				_, err := bot.Send(msg)
				if err != nil {
					log.Panic(err)
				}
				continue
			}

			// Make the API request

			botConfig := config.NewBotInfastructureConfig()

			resp, err := http.Get(fmt.Sprintf("https://holidays.abstractapi.com/v1/?api_key=13b720662b9e49f7a6c68d7e83f5e7ab&country=%s&year=%d", botConfig.APIHolidayKey, country, time.Now().Year()))
			if err != nil {
				logger.Errorf("Failed to make API request: %v", err)
				continue
			}

			// Read the response body
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Errorf("Failed to read response body: %v", err)
				continue
			}

			// Unmarshal the response body into an array of Holiday structs
			var holidays []Holiday
			err = json.Unmarshal(body, &holidays)
			if err != nil {
				logger.Errorf("Failed to unmarshal response body: %v", err)
				continue
			}

			// Construct the response message with the holidays
			var responseMsg strings.Builder
			responseMsg.WriteString("Upcoming holidays in ")
			responseMsg.WriteString(country)
			responseMsg.WriteString(":\n\n")

			for _, holiday := range holidays {
				// Parse the holiday date string into a time.Time object
				holidayDate, err := time.Parse("2006-01-02", holiday.Date)
				if err != nil {
					logger.Errorf("Failed to parse holiday date: %v", err)
					continue
				}

				// Format the holiday date as "Monday, January 2"
				holidayDateStr := holidayDate.Format("Monday, January 2")

				// Add the holiday name and date to the response message
				responseMsg.WriteString("- ")
				responseMsg.WriteString(holiday.Name)
				responseMsg.WriteString(" on ")
				responseMsg.WriteString(holidayDateStr)
				responseMsg.WriteString("\n")
			}

			// Send the response message
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, responseMsg.String())
			if _, err := bot.Send(msg); err != nil {
				logger.Errorf("Failed to send message: %v", err)
				continue
			}
		}
	}
}
