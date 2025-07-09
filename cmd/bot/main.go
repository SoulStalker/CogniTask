package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"

	tele "gopkg.in/telebot.v3"

	"github.com/SoulStalker/cognitask/internal/config"
	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/handlers"
	"github.com/SoulStalker/cognitask/internal/infra"
	"github.com/SoulStalker/cognitask/internal/middleware"
	"github.com/SoulStalker/cognitask/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm/logger"
)

func main() {
	// Контекст с отменой для graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//load config
	cfg := config.MustLoad()

	// init storage
	dbConfig := infra.DBConfig{DSN: cfg.DSN, Logger: logger.Default}
	db, err := infra.InitDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}

	// init repo
	taskRepo := infra.New(db)
	mediaRepo := infra.NewMediaRepo(db)
	settingRepo := infra.NewSettingsRepository(db)

	// init service
	taskUC := usecase.NewTaskService(taskRepo)
	mediaUC := usecase.NewMediaService(mediaRepo)
	settingsUC := usecase.NewSettingsService(settingRepo)

	// init redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.RedisHost, cfg.RedisPort),
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Проверка доступности редиски
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Successfully connected to Redis")

	// init fsm
	fsmService := fsm.NewFSMService(rdb, cfg.FSMTimeout)

	// init handler
	h := handlers.NewTaskHandler(fsmService, taskUC, ctx)
	mh := handlers.NewMediaHandler(mediaUC, ctx)
	sh := handlers.NewSettingsHandler(*settingsUC)
	cbRouter := handlers.NewCallbackRouter([]handlers.CallbackHandler{h}, fsmService, ctx)

	// run bot polling
	b, err := tele.NewBot(tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}
	if err = b.SetCommands(messages.Commands); err != nil {
		log.Fatalf("Не удалось установить команды: %v", err)
	}
	b.Use(middleware.AuthMiddleware(cfg.ChatId))

	// commands
	b.Handle("/start", h.Start)
	b.Handle("/help", h.Help)
	b.Handle("/random", mh.Random)

	// tasks
	b.Handle(tele.OnText, h.HandleText)

	b.Handle(keyboards.BtnAdd, h.Add)
	b.Handle(keyboards.BtnPending, h.Pending)
	b.Handle(keyboards.BtnComplete, h.Complete)
	b.Handle(keyboards.BtnDelete, h.Delete)
	b.Handle(keyboards.BtnAll, h.All)

	// media
	b.Handle(tele.OnMedia, mh.Create)
	// b.Handle(keyboards.BtnRandomPic, mh.Random)

	// settings
	b.Handle(keyboards.BtnSettings, sh.Settings)
	b.Handle(keyboards.BtnAutoDelete, sh.SetDeleteDays)

	// other
	b.Handle(tele.OnCallback, cbRouter.Handle)
	b.Handle(keyboards.BtnCancel, h.Cancel)

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		log.Println("Shutting down bot...")
		b.Stop()
		cancel()
	}()

	log.Println("Bot started")
	b.Start()
}
