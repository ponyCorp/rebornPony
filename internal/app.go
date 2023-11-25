package internal

import (
	"fmt"

	"github.com/ponyCorp/rebornPony/config"
	"github.com/ponyCorp/rebornPony/internal/repository"
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
	//!warn todo remove
	fmt.Println(conf)
	driver, err := repository.CreateDriver(conf.BdType)
	if err != nil {
		return err
	}
	err = driver.Connect(conf.BdPath, conf.DBName)
	if err != nil {
		return err
	}
	rep, err := repository.NewRepository(driver)
	if err != nil {
		return err
	}
	_ = rep
	return nil
}
