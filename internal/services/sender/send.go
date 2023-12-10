package sender

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Sender struct {
	Bot *tgbotapi.BotAPI
}

func NewSender(bot *tgbotapi.BotAPI) *Sender {
	return &Sender{
		Bot: bot,
	}
}
func (s *Sender) SendMessage(chatID int64, text string) error {
	_, err := s.Bot.Send(tgbotapi.NewMessage(chatID, text))
	return err
}

// Reply
func (s *Sender) Reply(chatId int64, msgID int, text string) error {
	msg := tgbotapi.NewMessage(chatId, text)
	msg.ReplyToMessageID = msgID
	_, err := s.Bot.Send(msg)
	return err
}
