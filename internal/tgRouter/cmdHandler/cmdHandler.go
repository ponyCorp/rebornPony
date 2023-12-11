package cmdhandler

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/internal/repository"
	"github.com/ponyCorp/rebornPony/internal/services/sender"
)

type handleFunc func(update *tgbotapi.Update, cmd, arg string)

type rout struct {
	handleFunc  handleFunc
	route       string
	description string
	groupName   string
}
type CmdHandler struct {
	mapRouter map[string]rout
	groups    map[string][]middlewareFunc
	rep       *repository.Repository
	rules     *rules
	sender    *sender.Sender
	// cmdParser *command.CommandParser
}

func NewCmdHandler(rep *repository.Repository, sender *sender.Sender) *CmdHandler {
	rules := newRules()
	rules.add("Руководство по командам бота")
	rules.add("Пользовательские команды:")
	return &CmdHandler{
		mapRouter: make(map[string]rout),
		groups:    make(map[string][]middlewareFunc),
		rep:       rep,
		rules:     rules,
		sender:    sender,
		//	cmdParser: command.NewCommandParser(botName),
	}
}

// add group middleware
func (h *CmdHandler) addGroupMiddleware(group string, f middlewareFunc) {
	if h.groups[group] == nil {
		h.groups[group] = make([]middlewareFunc, 0)
	}
	h.groups[group] = append(h.groups[group], f)
}
func (h *CmdHandler) handle(route string, description string, f handleFunc, group string) {

	h.rules.addRoute(route, description)
	h.mapRouter[route] = rout{
		handleFunc:  f,
		route:       route,
		description: description,
		groupName:   group,
	}
}
func (h *CmdHandler) Handle(route string, description string, f handleFunc) {
	h.handle(route, description, f, "")
}
func (h *CmdHandler) HadleUndefined(f handleFunc) {

	h.handle("undefined", "undefined", f, "")
}
func (h *CmdHandler) HandleMany(routes []string, description string, f handleFunc) {
	for _, route := range routes {
		h.Handle(route, description, f)
	}
}
func (h *CmdHandler) Route(update *tgbotapi.Update, cmd, arg string) bool {
	defer func() {
		if r := recover(); r != nil {
			fmc.Printfln("#fbtError>  #bbt[%+v]", r)
		}
	}()
	if rout, ok := h.mapRouter[cmd]; ok {
		if groupMiddlewares, ok := h.groups[rout.groupName]; ok {
			for _, groupMiddleware := range groupMiddlewares {
				ok, err := groupMiddleware(update, cmd, arg)
				if err != nil {
					fmc.Printfln("#fbtError>  #bbt[%+v]", err)
					return false
				}
				if !ok {
					return false
				}
			}
		}

		rout.handleFunc(update, cmd, arg)
		return true
	}
	if rout, ok := h.mapRouter["undefined"]; ok {
		rout.handleFunc(update, cmd, arg)
		return true
	}
	return false
}
