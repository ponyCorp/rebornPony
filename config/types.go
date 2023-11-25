package config

type Config struct {
	HostPort         string `env:"HOST_PORT"`
	DBHost           string `env:"DB_HOST"`
	DBName           string `env:"DB_NAME"`
	TelegramBotToken string `env:"TELEGRAM_BOT_TOKEN"`
	BdPath           string `env:"BD_PATH"`
	BdType           string `env:"BD_TYPE"`
	StoragePath      string `env:"STORAGE_PATH"`
}
