package main

import (
	"fmt"
	"github.com/TtMyth123/telegramBot"
	botClient "github.com/zelenin/grabot/client"
	"net/http"

	"github.com/astaxie/beego"
	"net/url"
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

func CallbackQueryFunc(methodName string, data map[string]interface{}) (string, error) {
	ss := ""
	for k, v := range data {
		ss += fmt.Sprintf("%s:%v\n", k, v)
	}
	str := fmt.Sprintf(`%s\n%s`, methodName, ss)
	return str, nil
}
