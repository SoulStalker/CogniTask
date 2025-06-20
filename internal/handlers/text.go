package handlers

import (
	"fmt"
	"log"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

// HandleText - это единый обработчик для всех текстовых сообщений.
// Он работает как маршрутизатор на основе состояния FSM.
// Оказывается telebot не умеет пропускать сообщения без обработки :(
func (h *TaskHandler) HandleText(c tele.Context) error {
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
		return c.Send("Unknown text: ", c.Text())
	}

}

func (h *TaskHandler) HandCallback(c tele.Context) error {

	fmt.Println("Unknown callback: ", c.Callback().Data)

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