package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const token = "*"
const groupID = int64(0)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		panic(err)
	}

	bot.Debug = true
	//u, e := bot.GetUpdates(tgbotapi.UpdateConfig{})
	msg := tgbotapi.NewMessage(groupID, "helo")
	u, e := bot.Send(msg)
	fmt.Println(u)
	fmt.Println(e)
}
