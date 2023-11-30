package cmdhandler

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type rules struct {
	msg strings.Builder
}

func newRules() *rules {
	return &rules{}
}
func (r *rules) add(msg string) {
	r.msg.WriteString(msg)
	r.msg.WriteString("\n")
}
func (r *rules) addRoute(route string, descrtiption string) {
	r.msg.WriteString(route)
	r.msg.WriteString(" - ")
	r.msg.WriteString(descrtiption)
	r.msg.WriteString("\n")

}
func (r *rules) get() string {
	return r.msg.String()
}
func (h *CmdHandler) Help(update *tgbotapi.Update, cmd, arg string) {
	h.sender.SendMessage(update.Message.Chat.ID, h.rules.get())
}
