package telegramBot

import (
	botClient "github.com/zelenin/grabot/client"
)

type CallbackQueryFunc func(methodName string, data map[string]interface{}) (botClient.AnswerCallbackQueryRequest, error)
type MessageFunc func(aMessage botClient.Message) error
