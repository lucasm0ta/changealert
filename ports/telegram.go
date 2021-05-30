package ports

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/lucasm0ta/trex-eyes/core"
)

//
type TelegramBot struct {
	//
	TelegramAccessToken string

	//
	watcherService *core.WatcherService
}

func buildWatchRequest(url string, chatId string) (*core.WatchRequest, error) {
	channel := core.NewChannelInfo("telegram", chatId)
	watchequest, err := core.NewWatchRequest(channel, url)
	if err != nil {
		return nil, err
	}
	return watchequest, nil
}

func NewTelegramBot(telegramAccessToken string, watcherService *core.WatcherService) *TelegramBot {
	telegramBot := new(TelegramBot)
	telegramBot.TelegramAccessToken = telegramAccessToken
	telegramBot.watcherService = watcherService
	return telegramBot
}

func (telegram *TelegramBot) Start() {
	bot, err := tgbotapi.NewBotAPI(telegram.TelegramAccessToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "watch":
				url := update.Message.CommandArguments()
				watchRequest, err := buildWatchRequest(url, fmt.Sprint(update.Message.Chat.ID))
				if err != nil {
					msg.Text = "oooops could not watch this"
					log.Panic(err)
				}
				err = telegram.watcherService.Register(watchRequest)
				if err == nil {
					msg.Text = "now watching"
				} else {
					msg.Text = "oooops could not watch this"
					log.Panic(err)

				}
			case "list ":
				msg.Text = "list"
			case "remove":
				msg.Text = "removed"
			case "help":
				msg.Text = "can't help U, sorry :/"
			default:
				msg.Text = ""
			}
		} else {
			msg.Text = `
			Sorry, I Con only understand the following commands:
			/watch
			/list
			/remove
			`
		}
		bot.Send(msg)
	}
}
