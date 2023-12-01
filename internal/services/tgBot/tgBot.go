package tgbot

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"
	"github.com/ponyCorp/rebornPony/internal/repository"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tgBot struct {
	rep   *repository.Repository
	token string
	Bot   *tgbotapi.BotAPI
	Upd   chan *tgbotapi.Update
	stop  chan bool
}

func NewTgBot(token string, rep *repository.Repository) *tgBot {

	return &tgBot{
		rep:   rep,
		token: token,
		Upd:   make(chan *tgbotapi.Update),
		stop:  make(chan bool),
	}
}
func (t *tgBot) Start() error {
	if t.token == "" {
		return errors.New("tgBot: token is empty")
	}
	fmt.Printf("tgBot: token: %s\n", t.token)
	//check regexp ^[0-9]{8,10}:[a-zA-Z0-9_-]{35}$
	match, err := regexp.MatchString(`^[0-9]{8,10}:[a-zA-Z0-9_-]{35}$`, t.token)
	if err != nil {
		return errors.Wrap(err, "regexp.MatchString")
	}
	if !match {
		return errors.New("tgBot: token is not valid")
	}
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		return errors.Wrap(err, "tgbotapi.NewBotAPI")
	}
	t.Bot = bot
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	go func() {
		for {
			select {
			case upd := <-updates:
				update := upd
				//fmt.Printf("tgBot: Update: %+v\n", upd)
				t.Upd <- &update
			//	fmt.Printf("tgBot: Update sended\n")
			case <-t.stop:
				fmt.Println("tgBot: Stop receiving updates")
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
func (t *tgBot) GetBotUsername() string {
	return t.Bot.Self.UserName
}
