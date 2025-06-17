package handlers

import (
	"github.com/SoulStalker/cognitask/internal/messages"
	"gopkg.in/telebot.v3"
)

func StartHandler(c telebot.Context) error {
	return c.Send(messages.BotMessages.Start)
}

func HelpHandler(c telebot.Context) error {
	return c.Send(messages.BotMessages.Help)
}

func AddHandler(c telebot.Context) error {
	return c.Send(messages.BotMessages.InputTaskText)
	// todo включаем состяние ожидания текста
}
