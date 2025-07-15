package infra

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DSN    string
	Logger logger.Interface
}

func InitDB(config DBConfig) (*gorm.DB, error) {
	var err error
	var db *gorm.DB
	db, err = gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: config.Logger,
	})
	if err != nil {
		log.Fatal(err)
	}

	// err = db.AutoMigrate(
	// 	&domain.Task{},
	// 	&domain.Media{},
	// 	&domain.Settings{}, 
	// )
	// if err != nil {
	// 	log.Fatal("AutoMigrate failed:", err)
	// }
	return db, err
}
