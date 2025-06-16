package main

import (
	"fmt"
	"log"
	"time"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/infra"
	"github.com/SoulStalker/cognitask/internal/usecase"
	"gorm.io/gorm/logger"
)

func main() {
	// init storage
	dbConfig := infra.DBConfig{DSN: "bot.db", Logger: logger.Default}
	db, err := infra.InitDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}


	// init repo
	taskRepo := infra.New(db)

	// init service
	taskService := usecase.NewTaskService(taskRepo)

	newTask, err := taskService.Add(domain.Task{
		Description: "Do something",
		Deadline: time.Now().AddDate(0, 0, 2),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(newTask)
	// init handler

	// run bot polling
}
