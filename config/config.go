package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/mallvielfrass/fmc"
	"github.com/ponyCorp/rebornPony/utils/file"
)

func lookup(envName string) (string, error) {
	env, exists := os.LookupEnv(envName)
	if !exists {
		return "", fmt.Errorf("env '%s' not found", envName)
	}
	return env, nil
}
func lookupInt(envName string) (int, error) {
	val, err := lookup(envName)
	if err != nil {
		return 0, err
	}
	parsedInt, err := strconv.Atoi(val)
	if err != nil {
		fmc.Printfln("#rbtError#wbt: #bbt%s", err.Error())
		return 0, err
	}
	return parsedInt, nil
}
func envLookup(envName string, defaultValue string) string {
	val, err := lookup(envName)
	if err != nil {
		fmc.Printfln("#rbtError#wbt: #bbt%s", err.Error())
		return defaultValue
	}
	return val

}
func loadEnvsFromFile(confPath string) error {
	fmc.Printfln("loadEnvsFromFile: %s", confPath)
	if confPath == "" {
		return nil
	}

	if !file.FileExists(confPath) {
		return fmt.Errorf("loadEnvsFromFile: no '%s' file not exist", confPath)
	}

	if err := godotenv.Load(confPath); err != nil {
		return fmt.Errorf("loadEnvsFromFile: no '%s' file open", confPath)
	}
	return nil
}
func InitConfig(confPath string) (Config, error) {
	err := loadEnvsFromFile(confPath)
	if err != nil {
		return Config{}, err
	}
	defaultConf := Config{
		HostPort: envLookup("HOST_PORT", ":9090"),
		DBHost:   envLookup("DB_HOST", "mongodb://127.0.0.1:27017/"),
		DBName:   envLookup("DB_NAME", "DefaultPonyDB"),

		TelegramBotToken: "",
		BdPath:           envLookup("BD_PATH", "../bd.db"),
		BdType:           envLookup("BD_TYPE", "sqlite"),
		StoragePath:      envLookup("STORAGE_PATH", "../"),
	}

	tgBotToken, err := lookup("TELEGRAM_BOT_TOKEN")
	if err != nil {
		return defaultConf, err
	}
	defaultConf.TelegramBotToken = tgBotToken

	return defaultConf, nil
}
