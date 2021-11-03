package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
	"github.com/tzY15368/lazarus/config"
)

var bot *tgbotapi.BotAPI
var err error
var debug = false

func ToggleBotDebug(v bool) {
	debug = v
}

func InitBot() error {
	bot, err = tgbotapi.NewBotAPI(config.Cfg.Master.TelegramAPIKey)
	if err != nil {
		return err
	}
	bot.Debug = debug
	return nil
}

func SendMessageToGroup(message string) {
	msg := tgbotapi.NewMessage(config.Cfg.Master.TelegramGroupID, message)
	_, e := bot.Send(msg)
	if e != nil {
		logrus.WithError(e).Warn("send message to group failed")
	}
}
