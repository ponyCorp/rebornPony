package tgrouter

import (
	"fmt"
	"log"
	"runtime/debug"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
)

type rout func(upd *tgbotapi.Update) error
type middlware func(upd *tgbotapi.Update, uType string) (bool, error)
type middlewareList struct {
	EventType string
	Func      middlware
}
type router struct {
	Routs       map[string]rout
	Middlewares []middlewareList
}

func NewRouter() router {
	return router{
		Routs:       make(map[string]rout),
		Middlewares: []middlewareList{},
	}
}
func (b *router) Handle(rout string, f rout) {
	b.Routs[rout] = f
}
func (b *router) Middleware(eventType string, f middlware) {
	b.Middlewares = append(b.Middlewares, middlewareList{
		EventType: eventType,
		Func:      f,
	})
}
func (b *router) middlewareIsExist(middle string) bool {
	for _, f := range b.Middlewares {
		if f.EventType == middle {
			return true
		}
	}
	return false
}
func (b *router) getMiddlewaresForFunc(localFuncName string) []middlewareList {
	var m []middlewareList
	for _, md := range b.Middlewares {
		if md.EventType == localFuncName || md.EventType == eventtypes.AllUpdateTypes {
			m = append(m, md)
		}
	}
	return m
}
func (b *router) getTypeUpdate(upd *tgbotapi.Update) string {
	return eventtypes.Undefined
}
func runMiddleware(m middlewareList, upd *tgbotapi.Update, localFuncName string) (b bool, returnErr error) {
	defer func() {
		if r := recover(); r != nil {
			b = false
			returnErr = r.(error)
		}
	}()

	ok, err := m.Func(upd, localFuncName)
	return ok, err
}
func (b *router) rangeAllMiddlewares(localFuncName string, upd *tgbotapi.Update) (bool, error) {
	middles := b.getMiddlewaresForFunc(localFuncName)
	for _, m := range middles {
		ok, err := runMiddleware(m, upd, localFuncName)
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}
func (b *router) stageMiddlewares(funcName string, update *tgbotapi.Update) (returnBool bool, returnErr error, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = r.(error)
			returnBool = false
			panicErr = r.(error)
		}
	}()
	res, rerr := b.rangeAllMiddlewares(funcName, update)
	return res, rerr, nil
}

func (b *router) stageRout(funcName string, update *tgbotapi.Update) (returnErr error, panicErr error) {
	defer func() {
		if r := recover(); r != nil {
			returnErr = r.(error)
			panicErr = r.(error)
		}
	}()
	f, ok := b.Routs[funcName]
	if !ok {
		return fmt.Errorf("#fbtUndefined function> #bbt[%+v]", funcName), nil
	}
	e := f(update)
	return e, nil
}
func (b *router) Run(upd *chan *tgbotapi.Update) error {
	//add default middleware
	if !b.middlewareIsExist(eventtypes.Undefined) {
		b.Middleware(eventtypes.Undefined, func(upd *tgbotapi.Update, uType string) (bool, error) {
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
