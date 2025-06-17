package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Bot_token            string        `yaml:"bot_token" env:"BOT_TOKEN" env-required:"true"`
	Chat_ID              int64         `yaml:"chat_id" env:"CHAT_ID" env-default:"123456789"`
	DSN                  string        `yaml:"dsn" env:"DATABASE_DSN" env-required:"true"`
	NotificationInterval time.Duration `yaml:"interval" env:"NOTIFICATION_INTERVAL_HOURS" env-default:"3h"`
	AutoDeleteDays       uint          `yaml:"auto_delete_days" env:"AUTO_DELETE_DAYS" env-default:"30"`
	NotificationStart    uint          `yaml:"notify_since" env:"NOTIFICATION_START" env-default:"9"`
	NotificationEnd      uint          `yaml:"notify_to" env:"NOTIFICATION_END" env-default:"19"`
	WebhookURL           string        `yaml:"webhook_server" env:"WEBHOOK_URL"`
	SeverPort            uint          `yaml:"server_port" env:"SERVER_PORT"`
	BotMode              string        `yaml:"bot_mode" env:"BOT_MODE" env-default:"polling"`
}

func MustLoad() *Config {
	configPath := "./config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
