package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tzY15368/lazarus/config"
)

var bot *tgbotapi.BotAPI
var err error

func InitBot() error {
	bot, err = tgbotapi.NewBotAPI(config.Cfg.Master.TelegramAPIKey)
	if err != nil {
		return err
	}
	bot.Debug = true
	return nil
}

func SendMessageToGroup(message string) error {
	msg := tgbotapi.NewMessage(config.Cfg.Master.TelegramGroupID, message)
	_, e := bot.Send(msg)
	if e != nil {
		return e
	}
	return nil
}
