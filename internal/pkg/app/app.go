package app

import (
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"hamsterbot/config"
	"hamsterbot/internal/app/handlers/base"
	"hamsterbot/internal/app/middleware"
	baseService "hamsterbot/internal/app/services/base"
	usersService "hamsterbot/internal/app/services/users"
	"hamsterbot/pkg/cache"
	"hamsterbot/pkg/db"
	"hamsterbot/pkg/logger"
	"log"
	"time"
)

type App struct {
	base  *baseService.Service
	users *usersService.Service
}

func New() (*App, error) {
	a := &App{}

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal("Ошибка при попытке спарсить .env файл в структуру", err.Error())
	}

	logger.Init(cfg.LoggerLevel)

	err = cache.Init(fmt.Sprintf("%s:%s", cfg.Redis.RedisAddr, cfg.Redis.RedisPort), cfg.Redis.RedisUsername, cfg.Redis.RedisPassword, cfg.Redis.RedisDBId)
	if err != nil {
		logger.Error("Ошибка при инициализации кэша", zap.Error(err))
	}

	err = db.Init(cfg.DB.DBUser, cfg.DB.DBPassword, cfg.DB.DBHost, cfg.DB.DBName)
	if err != nil {
		logger.Fatal("Ошибка при инициализации БД", zap.Error(err))
	}

	InitBot(cfg.TelegramAPI, a)

	return a, nil
}

func InitBot(TelegramAPI string, a *App) {
	pref := tele.Settings{
		Token:  TelegramAPI,
		Poller: &tele.LongPoller{Timeout: 1 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		logger.Fatal("Ошибка при создании бота", zap.Error(err), zap.Any("pref", pref))
	}

	// Сервисы
	a.base = baseService.New()
	a.users = usersService.New()

	// Middleware
	mw := middleware.Endpoint{Bot: b, User: a.users}
	b.Use(mw.IsUser)

	// Эндпоинты
	baseEndpoint := base.Endpoint{Base: a.base}

	// Обработчики
	b.Handle("/help", baseEndpoint.HelpHandler)

	logger.Info("Бот запущен")
	b.Start()
}
