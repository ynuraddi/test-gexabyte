package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

// TODO: дописать теги для перменный окружения
type Config struct {
	Server struct {
		Host string `env:"APP_HOST"`
		Port string `env:"APP_PORT"`
	}

	Postgres struct {
		DSN          string `env:"DB_DSN"`
		MigrationURL string `env:"MIGRATION_URL"`
	}

	LogLevel string `env:"LOG_LEVEL" env-default:"dev"`

	Binance struct {
		BaseURL   string `env:"BINANCE_BASE_URL"`
		ApiKey    string `env:"BINANCE_API_KEY"`
		SecretKey string `env:"BINANCE_SECRET_KEY"`
	}
}

// priority of configs: env > yaml > default
func MustLoad() *Config {
	var path string

	flag.StringVar(&path, "config_path", "", "path of config file")
	flag.Parse()

	if len(path) == 0 {
		path = os.Getenv("CONFIG_PATH")
	}

	var cfg Config

	if path != "" {
		if err := cleanenv.ReadConfig(path, &cfg); err != nil {
			panic("failed to read config: " + err.Error())
		}
	}

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}
