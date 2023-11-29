package internal

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/ponyCorp/rebornPony/config"
	"github.com/ponyCorp/rebornPony/internal/repository"
	tgbot "github.com/ponyCorp/rebornPony/internal/services/tgBot"
	tgrouter "github.com/ponyCorp/rebornPony/internal/tgRouter"
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
		return errors.Wrap(err, "initConfig")
	}

	driver, err := repository.CreateDriver(conf.BdType)
	if err != nil {
		return errors.Wrap(err, "createDriver")
	}

	err = driver.Connect(conf.BdPath, conf.DBName)
	if err != nil {
		return errors.Wrap(err, "connectDriver")
	}
	defer driver.Disconnect()
	rep, err := repository.NewRepository(driver)
	if err != nil {
		return errors.Wrap(err, "createRepository")
	}
	_ = rep
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)
	tgBot := tgbot.NewTgBot(conf.TelegramBotToken, rep)

	err = tgBot.Start()
	if err != nil {
		return errors.Wrap(err, "startTgBot")
	}
	tgRouter := tgrouter.NewRouter()
	tgRouter.Mount()
	fmt.Println("mounted")
	<-exit
	tgBot.Stop()
	fmt.Println("exit")
	return nil
}
