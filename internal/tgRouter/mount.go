package tgrouter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
)

func (r *Router) Mount() {
	r.Middleware(eventtypes.AllUpdateTypes, func(upd *tgbotapi.Update, uType string) (bool, error) {
		fmc.Printfln("#fbtAllUpdateTypes middleware> #bbt[%+v]", upd)
		return true, nil
	})
}
