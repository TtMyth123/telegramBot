package telegramBot

import (
	botClient "github.com/zelenin/grabot/client"
)

type CallbackQueryFunc func(aCallbackQuery botClient.CallbackQuery, data map[string]interface{}) (botClient.AnswerCallbackQueryRequest, error)
type MessageFunc func(aMessage botClient.Message) error
