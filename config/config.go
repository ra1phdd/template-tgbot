package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

type Configuration struct {
	TelegramAPI string `env:"TELEGRAM_API,required"`
	LoggerLevel string `env:"LOGGER_LEVEL" envDefault:"warn"`
	DB          DB
	Redis       Redis
}

type DB struct {
	Host     string `env:"DB_HOST,required"`
	Username string `env:"DB_USERNAME,required"`
	Password string `env:"DB_PASSWORD,required"`
	Name     string `env:"DB_NAME,required"`
}

type Redis struct {
	Address  string `env:"REDIS_ADDR,required"`
	Port     string `env:"REDIS_PORT" envDefault:"6379"`
	Username string `env:"REDIS_USERNAME,required"`
	Password string `env:"REDIS_PASSWORD,required"`
	DBId     int    `env:"REDIS_DB_ID,required"`
}

func NewConfig(files ...string) (*Configuration, error) {
	err := godotenv.Load(files...)
	if err != nil {
		log.Fatal("Файл .env не найден")
	}

	cfg := Configuration{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg.Redis)
	if err != nil {
		return nil, err
	}
	err = env.Parse(&cfg.DB)

	return &cfg, nil
}
