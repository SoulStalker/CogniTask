package main

import (
	"log"

	"github.com/SoulStalker/cognitask/internal/infra"
	"gorm.io/gorm/logger"
)

func main() {
	// init storage
	dbConfig := infra.DBConfig{DSN: "bot.db", Logger: logger.Default}
	_, err := infra.InitDB(dbConfig)
	if err != nil {
		log.Fatal(err)
	}


	// init repo

	// init service

	// init handler

	// run bot polling
}
