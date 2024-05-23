package telegram

import (
	"encoding/base64"
	"os"
	"sync"

	"github.com/Monska85/telegram-gateway/lib/logHelper"
	"github.com/Monska85/telegram-gateway/lib/messagesHandlers"
	"github.com/Monska85/telegram-gateway/lib/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot *tgbotapi.BotAPI
}

type routeNewMessage func(update tgbotapi.Update)

var instance *TelegramBot
var lock = &sync.Mutex{}
var logger = logHelper.GetInstance()

func (t TelegramBot) SendMessage(chatId int64, message string) (bool, error) {
	logger.Out(utils.LogInfo, "Sending message")
	msg := tgbotapi.NewMessage(chatId, message)
	_, err := t.bot.Send(msg)

	if err != nil {
		logger.Out(utils.LogError, err.Error())
		return false, err
	}
	return true, nil
}

func (t TelegramBot) SendImageMessage(chatId int64, message string, image string, imageName string) (bool, error) {
	logger.Out(utils.LogInfo, "Sending image message")

	// Decode base64 string to bytes
	imageBytes, err := base64.StdEncoding.DecodeString(image)
	if err != nil {
		logger.Out(utils.LogError, "Error decoding base64 image: "+err.Error())
		return false, err
	}

	// Set filename, use imageName if provided, otherwise default to "image.jpg"
	filename := "Image"
	if imageName != "" {
		filename = imageName
	}
	// Create a FileBytes object with the image bytes
	file := tgbotapi.FileBytes{
		Name:  filename,
		Bytes: imageBytes,
	}
	msg := tgbotapi.NewPhoto(chatId, file)
	msg.Caption = message
	_, err = t.bot.Send(msg)

	if err != nil {
		logger.Out(utils.LogError, err.Error())
		return false, err
	}
	return true, nil
}

func (t TelegramBot) ListenAndRoute(fn routeNewMessage) {
	logger.Out(utils.LogInfo, "Getting updates channel")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := t.bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			logger.Out(utils.LogInfo, "Routing new message to the router function")
			fn(update)
		}
	}
}

func GetInstance(token string) *TelegramBot {
	if instance != nil {
		logger.Out(utils.LogDebug, "Single instance of TelegramBot already created")
		return instance
	}

	lock.Lock()
	defer lock.Unlock()

	if instance != nil {
		logger.Out(utils.LogDebug, "Single instance of TelegramBot already created")
		return instance
	}

	logger.Out(utils.LogDebug, "Creating new Telegram bot instance")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		logger.Out(utils.LogError, "Error creating new Telegram bot")
		logger.Out(utils.LogError, err.Error())
		os.Exit(1)
	}

	logger.Out(utils.LogInfo, "Authorized on account", "bot", bot.Self.UserName)

	instance = &TelegramBot{bot: bot}

	// Start listening for messages
	go instance.ListenAndRoute(messagesHandlers.RouteNewMessage)

	return instance
}
