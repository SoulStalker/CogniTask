package handlers

import (
	"log"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

func (h *TaskHandler) HandCallback(c tele.Context) error {
	userID := c.Sender().ID

	state, err := h.fsmService.GetState(h.ctx, userID)
	if err != nil {
		log.Printf("Failed to get state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	switch state.State {
	case fsm.StateWaitingTaskText:
		return h.processTaskText(c, state)
	case fsm.StateWaitingTaskDate:
		return h.processTaskDate(c, state)
	default:
		return nil
	}

}
