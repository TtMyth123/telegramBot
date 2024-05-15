package main

import (
	"fmt"
	"github.com/TtMyth123/telegramBot"
	botClient "github.com/zelenin/grabot/client"
	"net/http"
	"net/url"
	"strconv"

	"github.com/astaxie/beego"
)

func main() {

	botToken := beego.AppConfig.String("Telegram::botToken")
	HttpProxyUrl := beego.AppConfig.String("Telegram::HttpProxyUrl")
	arrOption := make([]botClient.Option, 0)

	if HttpProxyUrl != "" {
		ProxyURL, _ := url.Parse(HttpProxyUrl)
		httpClient := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(ProxyURL),
			},
		}
		o := botClient.WithHttpClient(httpClient)
		arrOption = append(arrOption, o)
	}

	aBot, e := telegramBot.NewClient(botToken, arrOption...)
	if e != nil {
		fmt.Println(e)
	}

	aBot.MessageFunc = MessageFunc
	aBot.MCallbackQueryFunc = CallbackQueryFunc

	var phoneNumber string
	fmt.Scanln(&phoneNumber)
}

func MessageFunc(aMessage botClient.Message) error {
	fmt.Println(aMessage)
	return nil
}

func CallbackQueryFunc(methodName string, data map[string]interface{}) (botClient.AnswerCallbackQueryRequest, error) {
	ss := ""
	for k, v := range data {
		ss += fmt.Sprintf("%s:%v\n", k, v)
	}
	//str := fmt.Sprintf(`%s\n%s`, methodName, ss)
	req := botClient.AnswerCallbackQueryRequest{}
	ShowAlert := false
	req.ShowAlert = &ShowAlert
	//req.Text = &str
	aMessage := botClient.Message{}
	Text := "aaa\n\naa"
	aMessage.Text = &Text
	aMessage.Chat.Id = GetInterface2Int64(data[telegramBot.CallbackQueryKey_ChatId], 0)

	return req, nil
}

func GetInterface2Int64(mp interface{}, k int64) int64 {
	if r1, ok := mp.(int64); ok {
		return r1
	}
	if r1, ok := mp.(int); ok {
		return int64(r1)
	}
	if r1, ok := mp.(int32); ok {
		return int64(r1)
	}
	if r1, ok := mp.(float64); ok {
		return int64(r1)
	}

	if r1, ok := mp.(uint8); ok {
		return int64(r1)
	}
	if r1, ok := mp.(uint16); ok {
		return int64(r1)
	}
	if r1, ok := mp.(uint32); ok {
		return int64(r1)
	}
	if r1, ok := mp.(uint64); ok {
		return int64(r1)
	}
	if r1, ok := mp.(int8); ok {
		return int64(r1)
	}
	if r1, ok := mp.(int16); ok {
		return int64(r1)
	}
	if r1, ok := mp.(uint); ok {
		return int64(r1)
	}

	if r1, ok := mp.(string); ok {
		r, e := strconv.ParseInt(r1, 10, 64)
		if e == nil {
			return int64(r)
		} else {
			return k
		}
	}

	v := GetInterface2Str(mp, fmt.Sprint(k))
	r, e := strconv.ParseInt(v, 10, 64)
	if e == nil {
		return r
	} else {
		return k
	}

	return k
}

func GetInterface2Str(mp interface{}, k string) string {
	r := k
	if mp != nil {
		if r1, ok := mp.(string); ok {
			return r1
		} else {
			return fmt.Sprint(mp)
		}
	}

	return r
}
