package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// init service
	taskUC := usecase.NewTaskService(taskRepo)

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

	// run bot polling
	b, err := tele.NewBot(tele.Settings{
		Token:  cfg.BotToken,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}
	b.Use(middleware.AuthMiddleware(cfg.ChatId))
	b.Handle("/start", h.Start)
	b.Handle("/help", h.Help)
	b.Handle("/add", h.Add)
	b.Handle("/pending", h.Pending)
	b.Handle(tele.OnText, h.HandleText)


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
