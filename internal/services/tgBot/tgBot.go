package tgbot

import (
	"github.com/ponyCorp/rebornPony/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	rep   *repository.Repository
	token string
	Bot   *tgbotapi.BotAPI
	Upd   chan tgbotapi.Update
	stop  chan bool
}

func NewTgBot(token string, rep *repository.Repository) *tgBot {

	return &tgBot{
		rep:   rep,
		token: token,
		Upd:   make(chan tgbotapi.Update),
	}
}
func (t *tgBot) Start() error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return err
	}
	t.Bot = bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case upd := <-updates:
				t.Upd <- upd
			case <-t.stop:
				bot.StopReceivingUpdates()
				return
			}
		}
	}()
	return nil
}
func (t *tgBot) Stop() {
	t.stop <- true
}
