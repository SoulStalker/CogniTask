package handlers

import (
	"github.com/SoulStalker/cognitask/internal/messages"
	"gopkg.in/telebot.v3"
)

// сюда пишем код, который обрабатывает команды /add, /list и т.д.
func startHandler(c telebot.Context) error {
	return c.Send(messages.BotMessages.Start)
}
