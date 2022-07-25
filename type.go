package telegramBot

type Button struct {
	Text string
	Data string
	Url  string
}

const (
	CallbackQueryKey_UserId    = "UserId"
	CallbackQueryKey_ChatId    = "ChatId"
	CallbackQueryKey_Data      = "Data"
	CallbackQueryKey_MessageId = "MessageId"
)
