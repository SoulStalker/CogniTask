package handlers

import (
	"context"
	"log"
	"strconv"

	"github.com/SoulStalker/cognitask/internal/domain"
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

// Pending вывод списка открытых задач
func (h *TaskHandler) Pending(c tele.Context) error {
	tasks, err := h.service.GetPending()
	if err != nil {
		return c.Send(err.Error())
	}
	if len(tasks) == 0 {
		return c.Send("У вас нет открытых задач")
	}

	var rows [][]tele.InlineButton
	for _, task := range tasks {
		btn := 
		tele.InlineButton{
			Unique: "complete_task",
			Text:   "✅ " + task.Description,
			Data:   strconv.Itoa(int(task.ID)), // сюда попадёт в c.Callback().Data
		}
		rows = append(rows, []tele.InlineButton{btn})
	}

	return c.Send("Текущие задачи:", &tele.ReplyMarkup{InlineKeyboard: rows})
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
	default:
		return c.Send("Task number ", keyboards.CreateTaskKeyboard(1))
	}

}

// processTaskText в фсм состянии ждет название таска
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

// processTaskDate в фсм состянии ждет дату таски
func (h *TaskHandler) processTaskDate(c tele.Context, state *fsm.FSMData) error {
	userID := c.Sender().ID
	state.TaskDate = c.Text()

	taskDeadline, err := keyboards.ParseDate(state.TaskDate)
	if err != nil {
		return c.Send(err.Error())
	}

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to update state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	log.Printf("Task date saved: %s, moved to state: %s", state.TaskDate, state.State)
	log.Printf("Current state: %v", state)
	taskDescription := state.TaskText

	task, err := h.service.Add(domain.Task{
		Description: taskDescription,
		Deadline:    taskDeadline,
	})
	if err != nil {
		return c.Send(err.Error())
	}
	err = h.fsmService.ClearState(h.ctx, userID)
	if err != nil {
		return c.Send(err.Error())
	}
	return c.Send(formatTask(task))
}

// Complete закрытие задачи
func (h *TaskHandler) Complete(c tele.Context) error {
	strTaskID := c.Callback().Data
	taskID, err := strconv.Atoi(strTaskID)
	log.Printf("Task ID: %d", taskID)
	if err != nil {
		return c.Send(err.Error())
	}
	_, err = h.service.MarkDone(uint(taskID))
	if err != nil {
		return err
	}
	return c.Edit("Task completed")
}
