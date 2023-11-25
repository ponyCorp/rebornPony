package main

import (
	"flag"
	"fmt"

	"github.com/ponyCorp/rebornPony/internal"
)

func main() {
	useEnv := flag.Bool("env", false, "switch from env to config")
	configPath := flag.String("config", "./config.conf", "path to config file")
	flag.Parse()
	cPath := *configPath
	if *useEnv {
		fmt.Println("using env vars")
		cPath = ""
	}
	app := internal.NewApp(cPath)

	err := app.Run()
	if err != nil {
		panic(err)
	}
}
