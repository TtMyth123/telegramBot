package telegramBot

import (
	"fmt"
	"github.com/TtMyth123/telegramBot/kit"
	botClient "github.com/zelenin/grabot/client"
	"github.com/zelenin/grabot/updates"
	"log"
	"time"

	"context"
)

type TelegramBot struct {
	Client             *botClient.Client
	MCallbackQueryFunc CallbackQueryFunc
	MessageFunc        MessageFunc
}

func NewClient(token string, options ...botClient.Option) (*TelegramBot, error) {
	aTelegramBot := &TelegramBot{}
	Client, e := botClient.New(token, options...)
	if e != nil {
		return nil, e
	}
	aTelegramBot.Client = Client
	go aTelegramBot.run()

	return aTelegramBot, nil
}

func (this *TelegramBot) run() {
	ctx, _ := context.WithCancel(context.Background())
	longPoller := updates.NewLongPoller(this.Client)

	updatesChan, errsChan := longPoller.LongPoll(ctx, &botClient.GetUpdatesRequest{
		Offset: botClient.OptionalInt(0),
	}, 1*time.Second)

	for {
		select {
		case update := <-updatesChan:
			if update.CallbackQuery != nil {
				go func(aCallbackQuery1 botClient.CallbackQuery) {
					this.callbackQuery(aCallbackQuery1)
				}(*update.CallbackQuery)
			}
			if update.Message != nil {
				if this.MessageFunc != nil {
					this.MessageFunc(*update.Message)
				}
			}

		case err := <-errsChan:
			log.Printf("error: %s", err)
		}
	}
}

func (this *TelegramBot) callbackQuery(aCallbackQuery botClient.CallbackQuery) error {
	clickUserId := aCallbackQuery.From.Id
	if this.MCallbackQueryFunc != nil {
		ChatId := int64(0)
		MessageId := int64(0)
		if aCallbackQuery.Message != nil {
			ChatId = aCallbackQuery.Message.Chat.Id
			MessageId = aCallbackQuery.Message.MessageId
		}
		callbackData := ""
		if aCallbackQuery.Data != nil {
			callbackData = *aCallbackQuery.Data
		}

		mpData := make(map[string]interface{})
		mpData[CallbackQueryKey_UserId] = clickUserId
		mpData[CallbackQueryKey_ChatId] = ChatId
		mpData[CallbackQueryKey_Data] = callbackData
		mpData[CallbackQueryKey_MessageId] = MessageId

		req, e := this.MCallbackQueryFunc(aCallbackQuery, mpData)
		if e != nil {
			ShowAlert := true
			resultText := e.Error()
			req.Text = &resultText
			req.ShowAlert = &ShowAlert
			req.CallbackQueryId = aCallbackQuery.Id
		}
		if req.CallbackQueryId == "" {
			req.CallbackQueryId = aCallbackQuery.Id
		}

		//ShowAlert := true
		//req := &botClient.AnswerCallbackQueryRequest{}
		//req.Text = &resultText
		//req.ShowAlert = &ShowAlert
		//req.CallbackQueryId = aCallbackQuery.Id
		_, e = this.Client.AnswerCallbackQuery(&req)
		if e != nil {
			ShowAlert := true
			resultText := e.Error()
			req.Text = &resultText
			req.ShowAlert = &ShowAlert
			req.CallbackQueryId = aCallbackQuery.Id
			_, e := this.Client.AnswerCallbackQuery(&req)
			if e != nil {
				log.Println(e)
				return e
			}
		} else {
			return e
		}
	}
	return nil
}

func getReplyToMessageId(ReplyToMessageId int64) *int64 {
	if ReplyToMessageId == 0 {
		return nil
	}
	return &ReplyToMessageId
}
func (this *TelegramBot) SendMessage(ChatId, ReplyToMessageId int64, Text string, arrButton [][]Button) error {
	ReplyMarkup := this.getReplyMarkup(arrButton)

	req := &botClient.SendMessageRequest{
		ChatId:           botClient.IntChatId(ChatId),
		Text:             Text,
		ReplyToMessageId: getReplyToMessageId(ReplyToMessageId),
		ReplyMarkup:      ReplyMarkup,
	}
	_, e := this.Client.SendMessage(req)
	return e
}
func (this *TelegramBot) SendMessageContact(ChatId, ReplyToMessageId int64, PhoneNumber, FirstName string, arrButton [][]Button) error {
	ReplyMarkup := this.getReplyMarkup(arrButton)
	DisableNotification := true
	req := &botClient.SendContactRequest{
		ChatId:              botClient.IntChatId(ChatId),
		PhoneNumber:         PhoneNumber,
		FirstName:           FirstName,
		ReplyToMessageId:    getReplyToMessageId(ReplyToMessageId),
		ReplyMarkup:         ReplyMarkup,
		DisableNotification: &DisableNotification,
	}
	_, e := this.Client.SendContact(req)
	return e
}
func (this *TelegramBot) SendMessageFileLocalEx(ChatId int64, Path, Caption string, arrButton [][]Button) error {
	t := kit.GetFileType(Path)
	switch t {
	case kit.FT_mp3:
		return this.SendMessageAudioFileLocal(ChatId, Path, Caption, arrButton)
	case kit.FT_mp4:
		return this.SendMessageVideoFileLocal(ChatId, Path, Caption, arrButton)
	case kit.FT_gif:
		return this.SendMessageAnimationFileLocal(ChatId, Path, Caption, arrButton)
	case kit.FT_png:
		return this.SendMessagePhotoFileLocal(ChatId, Path, Caption, arrButton)

	}
	return fmt.Errorf("")
}
func (this *TelegramBot) SendMessagePhotoFileLocal(ChatId int64, Path, Caption string, arrButton [][]Button) error {
	var strCaption *string
	if Caption == "" {
		strCaption = &Caption
	}
	aFile, e := botClient.NewFileInputFile(Path)
	if e != nil {
		return e
	}

	ReplyMarkup := this.getReplyMarkup(arrButton)
	req := &botClient.SendPhotoRequest{
		ChatId:      botClient.IntChatId(ChatId),
		Caption:     strCaption,
		ReplyMarkup: ReplyMarkup,
		Photo:       aFile,
	}
	_, e = this.Client.SendPhoto(req)
	return e
}
func (this *TelegramBot) SendMessageAudioFileLocal(ChatId int64, Path, Caption string, arrButton [][]Button) error {
	var strCaption *string
	if Caption == "" {
		strCaption = &Caption
	}
	aFile, e := botClient.NewFileInputFile(Path)
	if e != nil {
		return e
	}

	ReplyMarkup := this.getReplyMarkup(arrButton)
	req := &botClient.SendAudioRequest{
		ChatId:      botClient.IntChatId(ChatId),
		Caption:     strCaption,
		ReplyMarkup: ReplyMarkup,
		Audio:       aFile,
	}
	_, e = this.Client.SendAudio(req)
	return e
}
func (this *TelegramBot) SendMessageVideoFileLocal(ChatId int64, Path, Caption string, arrButton [][]Button) error {
	var strCaption *string
	if Caption == "" {
		strCaption = &Caption
	}
	aFile, e := botClient.NewFileInputFile(Path)
	if e != nil {
		return e
	}

	ReplyMarkup := this.getReplyMarkup(arrButton)
	req := &botClient.SendVideoRequest{
		ChatId:      botClient.IntChatId(ChatId),
		Caption:     strCaption,
		ReplyMarkup: ReplyMarkup,
		Video:       aFile,
	}
	_, e = this.Client.SendVideo(req)
	return e
}
func (this *TelegramBot) SendMessageAnimationFileLocal(ChatId int64, Path, Caption string, arrButton [][]Button) error {
	var strCaption *string
	if Caption == "" {
		strCaption = &Caption
	}
	aFile, e := botClient.NewFileInputFile(Path)
	if e != nil {
		return e
	}

	ReplyMarkup := this.getReplyMarkup(arrButton)
	req := &botClient.SendAnimationRequest{
		ChatId:      botClient.IntChatId(ChatId),
		Caption:     strCaption,
		ReplyMarkup: ReplyMarkup,
		Animation:   aFile,
	}
	_, e = this.Client.SendAnimation(req)
	return e
}
func (this *TelegramBot) getReplyMarkup(arrButton [][]Button) *botClient.InlineKeyboardMarkup {
	if arrButton == nil {
		ReplyMarkup := &botClient.InlineKeyboardMarkup{}
		ReplyMarkup.InlineKeyboard = make([][]botClient.InlineKeyboardButton, 0)
		return ReplyMarkup
	}

	ReplyMarkup := &botClient.InlineKeyboardMarkup{}

	iLen := len(arrButton)
	ReplyMarkup.InlineKeyboard = make([][]botClient.InlineKeyboardButton, iLen)
	for i := 0; i < iLen; i++ {
		jLen := len(arrButton[i])
		ReplyMarkup.InlineKeyboard[i] = make([]botClient.InlineKeyboardButton, jLen)
		for j := 0; j < jLen; j++ {
			var CallbackData, Url *string

			if arrButton[i][j].Data != "" {
				CallbackData = &arrButton[i][j].Data
			}
			if arrButton[i][j].Url != "" {
				Url = &arrButton[i][j].Url
			}

			aInlineKeyboardButton := botClient.InlineKeyboardButton{}
			//aInlineKeyboardButton.SwitchInlineQueryCurrentChat
			aInlineKeyboardButton.CallbackData = CallbackData
			aInlineKeyboardButton.Url = Url
			aInlineKeyboardButton.Text = arrButton[i][j].Text
			ReplyMarkup.InlineKeyboard[i][j] = aInlineKeyboardButton
		}
	}

	return ReplyMarkup
}
