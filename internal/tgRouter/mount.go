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
	r.Handle(eventtypes.Message, func(update *tgbotapi.Update) error {
		fmc.Printf("#fbtMessage>#wbt%s > #gbt@%s > %s %s> #bbt[%+v]\n", update.FromChat().Title, update.Message.From.UserName, update.Message.From.FirstName, update.Message.From.LastName, update.Message.Text)
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
	cmdRouter.Handle("info", "info", func(update *tgbotapi.Update, cmd, arg string) {

		u, isAdmin, err := groupManager.GetAdminInfo(update.FromChat().ID, update.SentFrom().ID)
		if err != nil {
			if err.Error() == "record not found" {
				sender.Reply(update.FromChat().ID, update.Message.MessageID, "Ты не админ")
				return
			}
			fmc.Printfln("#fbtError>  #bbt[%+v]", err)
			return
		}
		if !isAdmin {
			sender.Reply(update.FromChat().ID, update.Message.MessageID, "Ты не админ")
			return
		}
		if u.Level == admintypes.Owner {
			sender.Reply(update.FromChat().ID, update.Message.MessageID, "Вы - владелец чата")
			return
		}

		sender.Reply(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Вы администратор %d уровня", u.Level))
	})
	cmdRouter.Handle("reloadowner", "reloadowner", func(update *tgbotapi.Update, cmd, arg string) {
		user, err := groupManager.GetOwner(update.FromChat().ID)
		if err != nil {
			fmc.Printfln("#fbtError>  #bbt[%+v]", err)
			return
		}
		sender.Reply(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Владелец бота: @%s", user.UserName))
	})

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
	// ownerGroup.Handle("levelup", "levelup", func(update *tgbotapi.Update, cmd, arg string) {
	// 	replyMsg := update.Message.ReplyToMessage
	// 	if replyMsg == nil {
	// 		sender.Reply(update.FromChat().ID, update.Message.MessageID, "Укажите сообщение")
	// 		return
	// 	}
	// 	targetUser := replyMsg.From.ID
	// 	userLevel, err := groupManager.GetUserStatus(update.FromChat().ID, targetUser)
	// 	if err != nil {
	// 		fmc.Printfln("#fbtError>  #bbt[%+v]", err)
	// 		return
	// 	}
	// 	if admintypes.MaxLevel() <= userLevel {
	// 		sender.Reply(update.FromChat().ID, update.Message.MessageID, "У пользователя максимальный уровень")
	// 		return
	// 	}
	// 	err = groupManager.SetAdminLevel(update.FromChat().ID, targetUser, userLevel.IncreaseLevel())
	// 	if err != nil {
	// 		fmc.Printfln("#fbtError>  #bbt[%+v]", err)
	// 		return
	// 	}
	// 	sender.Reply(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Пользователю %d повышен уровень привилегий", targetUser))

	// })
	// //leveldown
	// ownerGroup.Handle("leveldown", "leveldown", func(update *tgbotapi.Update, cmd, arg string) {
	// 	replyMsg := update.Message.ReplyToMessage
	// 	if replyMsg == nil {
	// 		sender.Reply(update.FromChat().ID, update.Message.MessageID, "Укажите сообщение")
	// 		return
	// 	}
	// 	targetUser := replyMsg.From.ID
	// 	userLevel, err := groupManager.GetUserStatus(update.FromChat().ID, targetUser)
	// 	if err != nil {
	// 		fmc.Printfln("#fbtError>  #bbt[%+v]", err)
	// 		return
	// 	}
	// 	if userLevel == admintypes.Owner {
	// 		sender.Reply(update.FromChat().ID, update.Message.MessageID, "нелья понизить владельца чата")
	// 		return
	// 	}
	// 	if userLevel == admintypes.User {
	// 		sender.Reply(update.FromChat().ID, update.Message.MessageID, "У пользователя минимальный уровень")
	// 		return
	// 	}
	// 	err = groupManager.SetAdminLevel(update.FromChat().ID, targetUser, userLevel.DecreaseLevel())
	// 	if err != nil {
	// 		fmc.Printfln("#fbtError>  #bbt[%+v]", err)
	// 		return
	// 	}
	// 	sender.Reply(update.FromChat().ID, update.Message.MessageID, fmt.Sprintf("Пользователю %d понижен уровень привилегий", targetUser))
	// })
	adminOnlyGroup := cmdRouter.NewGroup("admin")

	adminOnlyGroup.AddMiddleware(func(update *tgbotapi.Update, cmd, arg string) (bool, error) {
		fmc.Printfln("#fbtAdminOnly middleware> #bbt[%+v]", update)
		status, err := groupManager.GetUserStatus(update.FromChat().ID, update.SentFrom().ID)
		if err != nil {
			fmc.Printfln("#fbtError>  #bbt[%+v]", err)
			return false, err
		}
		if status == admintypes.User {
			sender.Reply(update.FromChat().ID, update.Message.MessageID, "У тебя недостаточный уровень привилегий")
			return false, nil
		}
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
		err = warnTrigger.AddWords(update.FromChat().ID, arg)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не добавлено", arg))
			return
		}
		sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("слово %s добавлено", arg))

	})
	//delwarn
	adminOnlyGroup.Handle("delwarn", "delwarn", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("delwarn")
		err := rep.ChatSensetive.RemoveChatSensetive(update.FromChat().ID, arg, sensetivetypes.Warn)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не удалено", arg))
			return
		}
		err = warnTrigger.DeleteWords(update.FromChat().ID, arg)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не удалено", arg))
			return
		}
		sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("слово %s удалено", arg))
	})
	adminOnlyGroup.Handle("addforbidden", "addforbidden", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("addforbidden")
		err := rep.ChatSensetive.AddChatSensetive(update.FromChat().ID, arg, sensetivetypes.Forbidden)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не добавлено", arg))
			return
		}
		err = forbiddenTrigger.AddWords(update.FromChat().ID, arg)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не добавлено", arg))
			return
		}
		sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("слово %s добавлено", arg))
	})
	adminOnlyGroup.Handle("delforbidden", "delforbidden", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("delforbidden")
		err := rep.ChatSensetive.RemoveChatSensetive(update.FromChat().ID, arg, sensetivetypes.Forbidden)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не удалено", arg))
			return
		}
		err = forbiddenTrigger.DeleteWords(update.FromChat().ID, arg)
		if err != nil {
			fmt.Println(err)
			sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("Произошла ошибка. слово %s не удалено", arg))
			return
		}
		sender.SendMessage(update.Message.Chat.ID, fmt.Sprintf("слово %s удалено", arg))
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
		list := rep.ChatSensetive.GetChatSensetiveWords(update.FromChat().ID, sensetivetypes.Forbidden)
		msg := ""
		for _, v := range list.Words {
			msg += v + "\n"
		}
		sender.SendMessage(update.Message.Chat.ID, msg)
	})
	adminOnlyGroup.Handle("enable", "enable", func(update *tgbotapi.Update, cmd, arg string) {
		fmt.Println("enable")
		sender.SendMessage(update.Message.Chat.ID, "Chat enabled")
	})

	return cmdRouter
}
