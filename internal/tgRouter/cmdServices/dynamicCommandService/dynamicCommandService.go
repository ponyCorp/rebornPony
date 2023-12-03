package dynamiccommandservice

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type DynamicCommandService struct {
}

func NewDynamicCommandService() *DynamicCommandService {
	return &DynamicCommandService{}
}
func (d *DynamicCommandService) Handle(update *tgbotapi.Update, cmd, arg string) bool {
	return false
}
