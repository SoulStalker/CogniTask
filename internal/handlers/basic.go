package handlers

import (
	"log"

	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

// Start хендлер для обработки команды Start
func (h *TaskHandler) Start(c tele.Context) error {
	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}
	return c.Send(messages.BotMessages.Start)
}

// Help хендлер для обработки команды Help
func (h *TaskHandler) Help(c tele.Context) error {
	return c.Send(messages.BotMessages.Help)
}

