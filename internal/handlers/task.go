package handlers

import (
	"context"
	"log"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/SoulStalker/cognitask/internal/usecase"
	tele "gopkg.in/telebot.v3"
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

// Add  хендлер для обработки команды  Add
func (h *TaskHandler) Add(c tele.Context) error {
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
	}

	return nil
}

// TaskName в фсм состянии ждет название таска
func (h *TaskHandler) processTaskText(c tele.Context, state *fsm.FSMData) error {
	userID := c.Sender().ID
	state.TaskText = c.Text()
	state.State = fsm.StateWaitingTaskDate

	// Надо сохранить новое состояние
	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to update state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	log.Printf("Task text saved: %s, moved to state: %s", state.TaskText, state.State)
	return c.Send(messages.BotMessages.InputNewDate, keyboards.GetDateSelectionKeyboard())
}

func (h *TaskHandler) processTaskDate(c tele.Context, state *fsm.FSMData) error {
	userID := c.Sender().ID
	state.TaskText = c.Text()
	state.State = fsm.StateWaitingTaskCategory

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to update state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	log.Printf("Task date saved: %s, moved to state: %s", state.TaskDate, state.State)
	log.Printf("Current state: !!!     %v    !!!", state)
	return c.Send(messages.BotMessages.InputNewDate, keyboards.GetDateSelectionKeyboard())

}
