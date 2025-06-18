package handlers

import (
	"context"
	"log"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/SoulStalker/cognitask/internal/usecase"
	"gopkg.in/telebot.v3"
)

type TaskHandler struct {
	fsmService *fsm.FSMService
	service    *usecase.TaskService
	ctx        context.Context
}

func NewTaskHandler(fsm *fsm.FSMService, service *usecase.TaskService, ctx context.Context) *TaskHandler {
	return &TaskHandler{
		fsmService: fsm,
		service:    service,
		ctx:        ctx,
	}
}

// Start хендлер для обработки команды Start
func (h *TaskHandler) Start(c telebot.Context) error {
	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}
	return c.Send(messages.BotMessages.Start)
}

// Help хендлер для обработки команды Help
func (h *TaskHandler) Help(c telebot.Context) error {
	return c.Send(messages.BotMessages.Help)
}

// Add  хендлер для обработки команды  Add
func (h *TaskHandler) Add(c telebot.Context) error {
	userID := c.Sender().ID

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	state := &fsm.FSMData{
		State: fsm.StateWaitingTaskText,
	}

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to set state: %v", err)
		c.Send(messages.BotMessages.ErrorSomeError)
	}
	return c.Send(messages.BotMessages.InputTaskText)
}

// TaskName в фсм состянии ждет название таска
func (h *TaskHandler) TaskName(c telebot.Context) error {
	// todo хорошо бы переделать под кнопки с датами
	userID := c.Sender().ID

	// Текущще состоние фсм
	state, err := h.fsmService.GetState(h.ctx, userID)
	if err != nil {
		log.Printf("Failed to get state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	// Сверяем состояние
	if state.State == fsm.StateWaitingTaskText {
		state.TaskText = c.Text()
		// переходим на следующее состояние
		state.State = fsm.StateWaitingTaskDate

		// Надо сохранить новое состояние
		if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
			log.Printf("Failed to update state: %v", err)
			return c.Send(messages.BotMessages.ErrorTryAgain)
		}

		log.Printf("Task text saved: %s, moved to state: %s", state.TaskText, state.State)
		return c.Send(messages.BotMessages.InputNewDate)
	}

	// Если пользователь не в состоянии FSM, игнорируем сообщение
	// Или можно обработать как обычный текст
	return nil
}
