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
	printTask(newTask)
	newTask2, err := taskService.Add(domain.Task{
		Description: "Do something else",
		Deadline: time.Now().AddDate(0, 0, 1),
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	printTask(newTask2)
	_, err = taskService.MarkDone(newTask.ID)
	if err != nil {
		fmt.Println(err)
	}
	tasks, _ := taskService.GetPending()
	for _, task := range tasks {
		printTask(task)
	}
	// init handler

	// run bot polling


	
}

func printTask(task domain.Task) {
		fmt.Printf("Task: \nID: %d\nDesc: %s\nDeadline: %v\nCreatedAt: %v\nClosed: %v\nClosedAt: %v\n", task.ID, task.Description, task.Deadline, task.CreatedAt, task.Closed, task.ClosedAt)
	}