package main

import (
	"flag"

	"github.com/ponyCorp/rebornPony/internal"
)

func main() {
	useEnv := flag.Bool("env", false, "switch from env to config")
	configPath := flag.String("config", "./config.conf", "path to config file")
	cPath := *configPath
	if *useEnv {
		cPath = ""
	}
	app := internal.NewApp(cPath)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
