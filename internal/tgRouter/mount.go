package tgrouter

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/models/admintypes"
	"github.com/ponyCorp/rebornPony/internal/models/sensetivetypes"
	"github.com/ponyCorp/rebornPony/internal/repository"
	adminsmanager "github.com/ponyCorp/rebornPony/internal/services/adminsManager"
	eventchatswitcher "github.com/ponyCorp/rebornPony/internal/services/eventChatSwitcher"
	"github.com/ponyCorp/rebornPony/internal/services/sender"
	cmdhandler "github.com/ponyCorp/rebornPony/internal/tgRouter/cmdHandler"
	dynamiccommandservice "github.com/ponyCorp/rebornPony/internal/tgRouter/cmdServices/dynamicCommandService"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
	sensetivetrigger "github.com/ponyCorp/rebornPony/internal/tgRouter/sensetiveTrigger"
	"github.com/ponyCorp/rebornPony/utils/command"
)

func (r *Router) Mount(rep *repository.Repository, switcher *eventchatswitcher.EventChatSwitcher, sender *sender.Sender, groupManager *adminsmanager.AdminsManager, botName string) {
	cmdParser := command.NewCommandParser(botName)
	sensWarn := rep.ChatSensetive.GetAllChatSensetiveWords(sensetivetypes.Warn)
	warnTrigger := sensetivetrigger.New()
	for chat, v := range sensWarn {
		//	fmc.Printfln("#fbtWarn> CHAT[%+v] #bbt[%+v]", chat, v)
		err := warnTrigger.AddWords(chat, v...)
		if err != nil {
			fmc.Printfln("#rbtError>  #bbt[%+v]", err)
		}
	}

	sensForbidden := rep.ChatSensetive.GetAllChatSensetiveWords(sensetivetypes.Forbidden)
	forbiddenTrigger := sensetivetrigger.New()
	for chat, v := range sensForbidden {
		forbiddenTrigger.AddWords(chat, v...)
	}
	cmdRouter := r.cmdRouts(rep, sender, warnTrigger, forbiddenTrigger, groupManager)

	r.Handle(eventtypes.CommandMessage, func(update *tgbotapi.Update) error {
		isCommand, cmd := cmdParser.ParseCommand(update.Message.Text)
		if !isCommand {
			return nil
		}
		cmdRouter.Route(update, cmd.Cmd, cmd.Arg)
		return nil
	})
	r.Middleware(eventtypes.Message, func(upd *tgbotapi.Update, event eventtypes.Event) (bool, error) {
		fmc.Printf("#fbtMessage middleware> #bbt[%+v]\n", upd)
		//sensetive warn
		isSensetiveWarn := warnTrigger.ChatIsSensetive(upd.FromChat().ID)
		// fmt.Printf("isSensetiveWarn: %v\n", isSensetiveWarn)
		if !isSensetiveWarn {
			return true, nil
		}
		sesnD, err := warnTrigger.MessageContainSensitiveWords(upd.FromChat().ID, upd.Message.Text)
		if err != nil {
			return false, err
		}

		if sesnD {
			sender.Reply(upd.FromChat().ID, upd.Message.MessageID, "Ваше сообщение содержит нежелательные слова")
			return false, nil
		}

		//sensetive forbidden
		isSensetiveForbidden := forbiddenTrigger.ChatIsSensetive(upd.FromChat().ID)
		if !isSensetiveForbidden {
			return true, nil
		}
		sensF, err := forbiddenTrigger.MessageContainSensitiveWords(upd.FromChat().ID, upd.Message.Text)
		if err != nil {
			return false, err
		}

		if sensF {
			sender.Reply(upd.FromChat().ID, upd.Message.MessageID, "Ваше сообщение содержит запрещенные слова")
			return false, nil
		}
		return true, nil
	})

	r.Middleware(eventtypes.AllUpdateTypes, func(upd *tgbotapi.Update, event eventtypes.Event) (bool, error) {
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
func (r *Router) cmdRouts(rep *repository.Repository, sender *sender.Sender, warnTrigger *sensetivetrigger.SensetiveTrigger, forbiddenTrigger *sensetivetrigger.SensetiveTrigger, groupManager *adminsmanager.AdminsManager) *cmdhandler.CmdHandler {
	cmdRouter := cmdhandler.NewCmdHandler(rep, sender)
	dynamicService := dynamiccommandservice.NewDynamicCommandService()
	cmdRouter.Handle("help", "help", cmdRouter.Help)
	cmdRouter.HadleUndefined(func(update *tgbotapi.Update, cmd, arg string) {
		//sender.SendMessage(update.Message.Chat.ID, "Unknown command")
		dynamicService.Handle(update, cmd, arg)
	})
	ownerGroup := cmdRouter.NewGroup("owner")
	ownerGroup.AddMiddleware(func(update *tgbotapi.Update, cmd, arg string) (bool, error) {
		fmc.Printfln("#gbtOwnerMiddleware>  #bbt[%+v]", update)
		uLevel, err := groupManager.GetUserStatus(update.FromChat().ID, update.SentFrom().ID)
		if err != nil {
			fmc.Printfln("#fbtError>  #bbt[%+v]", err)
			return false, err
		}
		if uLevel < admintypes.Owner {
			sender.Reply(update.FromChat().ID, update.Message.MessageID, "У тебя недостаточный уровень привилегий")
			return false, nil
		}
		return true, nil
	})
	adminOnlyGroup := cmdRouter.NewGroup("admin")

	adminOnlyGroup.AddMiddleware(func(update *tgbotapi.Update, cmd, arg string) (bool, error) {
		fmc.Printfln("#fbtAdminOnly middleware> #bbt[%+v]", update)

		return true, nil
	})
	adminOnlyGroup.Handle("addwarn", "addwarn", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("addwarn")
		err := rep.ChatSensetive.AddChatSensetive(update.FromChat().ID, arg, sensetivetypes.Warn)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не добавлено", arg))
			return
		}
		warnTrigger.AddWords(update.FromChat().ID, arg)
		sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("слово %s добавлено", arg))

	})
	adminOnlyGroup.Handle("addforbidden", "addforbidden", func(update *tgbotapi.Update, cmd, arg string) {

	})
	//listwarn
	adminOnlyGroup.Handle("listwarn", "listwarn", func(update *tgbotapi.Update, cmd, arg string) {
		list := rep.ChatSensetive.GetChatSensetiveWords(update.FromChat().ID, sensetivetypes.Warn)
		msg := ""
		for _, v := range list.Words {
			msg += v + "\n"
		}
		sender.SendMessage(update.Message.Chat.ID, msg)
	})
	adminOnlyGroup.Handle("listforbidden", "listforbidden", func(update *tgbotapi.Update, cmd, arg string) {
	})
	adminOnlyGroup.Handle("enable", "enable", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("enable")
		sender.SendMessage(update.Message.Chat.ID, "Chat enabled")
	})

	return cmdRouter
}
