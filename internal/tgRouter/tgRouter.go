package tgrouter

import (
	"fmt"
	"log"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
	"github.com/ponyCorp/rebornPony/utils/command"
)

type rout func(upd *tgbotapi.Update) error
type middlware func(upd *tgbotapi.Update, event eventtypes.Event) (bool, error)
type middlewareList struct {
	EventType string
	Func      middlware
}
type Router struct {
	Routs       map[string]rout
	Middlewares []middlewareList
	cmdParser   *command.CommandParser
}

func NewRouter(botName string) Router {
	return Router{
		Routs:       make(map[string]rout),
		Middlewares: []middlewareList{},
		cmdParser:   command.NewCommandParser(botName),
	}
}
func (b *Router) Handle(event eventtypes.Event, f rout) {
	b.Routs[event.String()] = f
}
func (b *Router) Middleware(event eventtypes.Event, f middlware) {
	b.Middlewares = append(b.Middlewares, middlewareList{
		EventType: event.String(),
		Func:      f,
	})
}
func (b *Router) middlewareIsExist(event eventtypes.Event) bool {
	for _, f := range b.Middlewares {
		if f.EventType == event.String() {
			return true
		}
	}
	return false
}
func (b *Router) getMiddlewaresForFunc(event eventtypes.Event) []middlewareList {
	var m []middlewareList
	for _, md := range b.Middlewares {
		if md.EventType == event.String() || md.EventType == eventtypes.AllUpdateTypes.String() {
			m = append(m, md)
		}
	}
	return m
}
func (b *Router) getTypeUpdate(upd *tgbotapi.Update) eventtypes.Event {

	if upd.Message != nil {
		if upd.Message.NewChatMembers != nil {
			return eventtypes.MemberJoinInChat
		}
		if upd.Message.LeftChatMember != nil {
			return eventtypes.LeaveChatMember
		}
		if upd.Message.ForwardFrom != nil {
			return eventtypes.Message
		}

		if b.cmdParser.IsCommand(upd) {
			return eventtypes.CommandMessage
		}

		if upd.Message.ReplyToMessage != nil {
			if _, ok := b.Routs[eventtypes.ReplyMessage.String()]; ok {
				return eventtypes.ReplyMessage
			}
		}
		return eventtypes.Message
	}

	if upd.CallbackQuery != nil {
		return eventtypes.CallbackQuery
	}

	return eventtypes.Undefined
}
func runMiddleware(m middlewareList, upd *tgbotapi.Update, event eventtypes.Event) (b bool, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			b = false
			returnErr = r.(error)
		}
	}()

	ok, err := m.Func(upd, event)
	return ok, err
}
func (b *Router) rangeAllMiddlewares(event eventtypes.Event, upd *tgbotapi.Update) (bool, error) {
	middles := b.getMiddlewaresForFunc(event)
	for _, m := range middles {
		ok, err := runMiddleware(m, upd, event)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
func (b *Router) stageMiddlewares(event eventtypes.Event, update *tgbotapi.Update) (returnBool bool, returnErr error, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = r.(error)
			returnBool = false
			panicErr = r.(error)
		}
	}()
	res, rerr := b.rangeAllMiddlewares(event, update)
	return res, rerr, nil
}

func (b *Router) stageRout(event eventtypes.Event, update *tgbotapi.Update) (returnErr error, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = r.(error)
			panicErr = r.(error)
		}
	}()
	f, ok := b.Routs[event.String()]
	if !ok {
		return fmt.Errorf("#fbtUndefined function> #bbt[%+v]", event.String()), nil
	}
	e := f(update)
	return e, nil
}
func (b *Router) Run(upd *chan *tgbotapi.Update) error {
	//add default middleware
	if !b.middlewareIsExist(eventtypes.Undefined) {
		b.Middleware(eventtypes.Undefined, func(upd *tgbotapi.Update, event eventtypes.Event) (bool, error) {
			fmc.Printfln("#fbtUndefined middleware> #bbt[%+v]", upd)
			return false, nil
		})
	}
	for update := range *upd {
		funcName := b.getTypeUpdate(update)
		//run middlewares
		ok, err, panics := b.stageMiddlewares(funcName, update)
		if panics != nil {
			log.Println("stageMiddlewares| panic handle occurred:", panics)
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			continue
		}
		if err != nil {
			fmc.Printfln("stageMiddlewares| Error: [%+v]", err)
			continue
		}
		if !ok {
			continue
		}
		//run rout function
		err, panicErr := b.stageRout(funcName, update)
		if panicErr != nil {
			log.Println("stageRout| panic handle occurred:", panicErr)
			fmt.Println("stacktrace from panic: \n" + string(debug.Stack()))
			continue
		}
		if err != nil {
			fmc.Printfln("stageRout| Error: [%+v]", err)
			continue
		}
	}
	return nil
}
