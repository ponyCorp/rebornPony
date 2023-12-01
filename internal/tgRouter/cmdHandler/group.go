package cmdhandler

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type middlewareFunc func(update *tgbotapi.Update, cmd, arg string) (bool, error)

type group struct {
	groupName   string
	middlewares []middlewareFunc
	h           *CmdHandler
}

func (h *CmdHandler) NewGroup(groupName string) *group {
	return &group{
		groupName:   groupName,
		middlewares: make([]middlewareFunc, 0),
		h:           h,
	}
}
func (g *group) AddMiddleware(f middlewareFunc) {
	g.middlewares = append(g.middlewares, f)
}
func (g *group) Handle(route string, description string, f handleFunc) {
	g.h.handle(route, description, f, g.groupName)
}
func (g *group) HandleMany(routes []string, description string, f handleFunc) {
	for _, route := range routes {
		g.Handle(route, description, f)
	}
}
