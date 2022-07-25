package telegramBot

import (
	botClient "github.com/zelenin/grabot/client"
)

type CallbackQueryFunc func(methodName string, data map[string]interface{}) (string, error)
type MessageFunc func(aMessage botClient.Message) error
