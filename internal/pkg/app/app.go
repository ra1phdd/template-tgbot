package app

import (
	"fmt"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"hamsterbot/config"
	"hamsterbot/internal/app/handlers/base"
	baseService "hamsterbot/internal/app/services/base"
	"hamsterbot/pkg/cache"
	"hamsterbot/pkg/db"
	"hamsterbot/pkg/logger"
	"log"
	"strings"
	"time"
)

type App struct {
	base *baseService.Service
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

	// Эндпоинты
	baseEndpoint := base.Endpoint{Base: a.base}

	// Middlewares
	b.Use(mwUsers.IsUser)

	// Обработчики
	b.Handle("/help", baseEndpoint.HelpHandler)

	b.Handle("/send", func(c tele.Context) error {
		if c.Sender().ID != 1230045591 {
			return nil
		}

		args := c.Args()

		chatID := int64(-1002138316635)

		// Используем метод Send у объекта бота для отправки сообщения
		_, err := c.Bot().Send(tele.ChatID(chatID), strings.Join(args, " "))
		return err
	})

	b.Handle(tele.OnText, func(c tele.Context) error { return nil })
	b.Handle(tele.OnAudio, func(c tele.Context) error { return nil })
	b.Handle(tele.OnCallback, func(c tele.Context) error { return nil })
	b.Handle(tele.OnDocument, func(c tele.Context) error { return nil })
	b.Handle(tele.OnEdited, func(c tele.Context) error { return nil })
	b.Handle(tele.OnMedia, func(c tele.Context) error { return nil })
	b.Handle(tele.OnPhoto, func(c tele.Context) error { return nil })
	b.Handle(tele.OnSticker, func(c tele.Context) error { return nil })
	b.Handle(tele.OnVideo, func(c tele.Context) error { return nil })
	b.Handle(tele.OnVoice, func(c tele.Context) error { return nil })

	b.Start()
}
