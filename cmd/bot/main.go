package main

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v3"

	"github.com/SoulStalker/cognitask/internal/config"
	"github.com/SoulStalker/cognitask/internal/infra"
	"github.com/SoulStalker/cognitask/internal/usecase"
	"gorm.io/gorm/logger"
)

func main() {
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
	_ = usecase.NewTaskService(taskRepo)

	// init handler

	// run bot polling

	b, err := tele.NewBot(tele.Settings{
		Token:  cfg.Bot_token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return
	}
	b.Handle(tele.OnText, func(ctx tele.Context) error {
		return ctx.Send("Hello")
	})
	b.Start()

}
