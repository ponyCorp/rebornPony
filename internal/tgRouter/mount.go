package tgrouter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/repository"
	eventchatswitcher "github.com/ponyCorp/rebornPony/internal/services/eventChatSwitcher"
	cmdhandler "github.com/ponyCorp/rebornPony/internal/tgRouter/cmdHandler"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
	"github.com/ponyCorp/rebornPony/utils/command"
)

func (r *Router) Mount(rep *repository.Repository, switcher *eventchatswitcher.EventChatSwitcher, botName string) {
	cmdParser := command.NewCommandParser(botName)
	cmdRouter := cmdhandler.NewCmdHandler(rep)
	r.Handle(eventtypes.CommandMessage, func(update *tgbotapi.Update) error {
		isCommand, cmd := cmdParser.ParseCommand(update.Message.Text)
		if !isCommand {
			return nil
		}
		cmdRouter.Route(update, cmd.Cmd, cmd.Arg)
		return nil
	})
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
