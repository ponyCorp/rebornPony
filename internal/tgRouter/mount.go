package tgrouter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/repository"
	eventchatswitcher "github.com/ponyCorp/rebornPony/internal/services/eventChatSwitcher"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
)

func (r *Router) Mount(rep *repository.Repository, switcher *eventchatswitcher.EventChatSwitcher) {
	r.Middleware(eventtypes.AllUpdateTypes, func(upd *tgbotapi.Update, uType string) (bool, error) {
		fmc.Printfln("#fbtAllUpdateTypes middleware> #bbt[%+v]", upd)
		return true, nil
	})
	r.Middleware(eventtypes.AllUpdateTypes, func(upd *tgbotapi.Update, uType string) (bool, error) {
		isDisabled, err := switcher.ChatIsDisabled(upd.FromChat().ID)
		if err != nil {
			return false, err
		}

		if isDisabled {
			if upd.Message != nil && upd.Message.IsCommand() && upd.Message.Command() == "enable" {
				return true, nil
			}

			return false, nil
		}
		return true, nil
	})
}
