package tgrouter

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/repository"
	eventchatswitcher "github.com/ponyCorp/rebornPony/internal/services/eventChatSwitcher"
	"github.com/ponyCorp/rebornPony/internal/services/sender"
	cmdhandler "github.com/ponyCorp/rebornPony/internal/tgRouter/cmdHandler"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
	"github.com/ponyCorp/rebornPony/utils/command"
)

func (r *Router) Mount(rep *repository.Repository, switcher *eventchatswitcher.EventChatSwitcher, sender *sender.Sender, botName string) {
	cmdParser := command.NewCommandParser(botName)
	cmdRouter := r.cmdRouts(rep, sender)
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
		if !isDisabled {
			return true, nil
		}

		_, args := r.cmdParser.IsCommandWithArgs(upd)
		if args.Cmd == "enable" {
			return true, nil
		}
		return false, nil

	})
}
func (r *Router) cmdRouts(rep *repository.Repository, sender *sender.Sender) *cmdhandler.CmdHandler {
	cmdRouter := cmdhandler.NewCmdHandler(rep, sender)
	cmdRouter.Handle("help", "help", cmdRouter.Help)
	adminOnlyGroup := cmdRouter.NewGroup("admin")
	adminOnlyGroup.AddMiddleware(func(update *tgbotapi.Update, cmd, arg string) (bool, error) {
		fmc.Printfln("#fbtAdminOnly middleware> #bbt[%+v]", update)
		return true, nil
	})
	adminOnlyGroup.Handle("enable", "enable", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("enable")
		sender.SendMessage(update.Message.Chat.ID, "Chat enabled")
	})

	return cmdRouter
}
