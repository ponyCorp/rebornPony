package internal

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/config"
)

type App struct {
	confPath string
}

func NewApp(confPath string) *App {
	return &App{}
}
func (app *App) Run() error {
	conf, err := config.InitConfig(app.confPath)
	if err != nil {
		return err
	}
	fmt.Println(conf)
	driver, err := repository.CreateDriver(conf.BdType)
	if err != nil {
		return err
	}
	err = driver.Connect(conf.BdPath, conf.DBName)
	if err != nil {
		return err
	}
	return nil
}
