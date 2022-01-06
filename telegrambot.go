package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func sendMsg(msg string, bot *tgbotapi.BotAPI) {
	fullMsg := tgbotapi.NewMessageToChannel("@ashishxbox", msg)
	resmsg, err := bot.Send(fullMsg)
	if err != nil {
		fmt.Printf("err sending %#v msg %#v", resmsg.MessageID, err)
	}
}
func getBot() (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI("telegram api key")
	if err != nil {
		fmt.Println("error in getting bot")
		return nil, err
	}
	return bot, nil
}
