package internal

import (
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/config"
	"github.com/ponyCorp/rebornPony/internal/repository"
	tgbot "github.com/ponyCorp/rebornPony/internal/services/tgBot"
	tgrouter "github.com/ponyCorp/rebornPony/internal/tgRouter"
	eventtypes "github.com/ponyCorp/rebornPony/internal/tgRouter/eventTypes"
)

type App struct {
	confPath string
}

func NewApp(confPath string) *App {
	return &App{
		confPath: confPath,
	}
}
func (app *App) Run() error {
	conf, err := config.InitConfig(app.confPath)
	if err != nil {
		return err
	}

	driver, err := repository.CreateDriver(conf.BdType)
	if err != nil {
		return err
	}

	err = driver.Connect(conf.BdPath, conf.DBName)
	if err != nil {
		return err
	}
	defer driver.Disconnect()
	rep, err := repository.NewRepository(driver)
	if err != nil {
		return err
	}
	_ = rep
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	tgBot := tgbot.NewTgBot(conf.TelegramBotToken, rep)

	err = tgBot.Start()
	if err != nil {
		return err
	}
	tgRouter := tgrouter.NewRouter()
	tgRouter.Middleware(eventtypes.AllUpdateTypes, func(upd *tgbotapi.Update, uType string) (bool, error) {
		fmc.Printfln("#fbtAllUpdateTypes middleware> #bbt[%+v]", upd)
		return true, nil
	})

	<-exit
	tgBot.Stop()

	return nil
}
