package infra

import (
	"github.com/SoulStalker/cognitask/internal/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DSN    string
	Logger logger.Interface
}

func InitDB(config DBConfig) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(config.DSN), &gorm.Config{
		Logger: config.Logger,
	})

	db.AutoMigrate(&domain.Task{}, &domain.Media{}, domain.Settings{})
	return db, err
}
