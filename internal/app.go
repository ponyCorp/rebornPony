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
	return nil
}
