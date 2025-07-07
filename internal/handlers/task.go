package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
	err := c.Respond()
	if err != nil {
		return err
	}

	tasks, err := h.service.GetPending()
	if err != nil {
		return c.Send(err.Error())
	}
	if len(tasks) == 0 {
		return c.Send("У вас нет открытых задач")
	}

	rows := formatTaskList(tasks)

	return c.Edit("Текущие задачи:", &tele.ReplyMarkup{InlineKeyboard: rows})
}

// Add  хендлер для обработки команды  Add
func (h *TaskHandler) Add(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	userID := c.Sender().ID

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	state := &fsm.FSMData{
		State: fsm.StateWaitingTaskText,
	}

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to set state: %v", err)
		c.Edit(messages.BotMessages.ErrorSomeError)
	}
	return c.Edit(messages.BotMessages.InputTaskText, keyboards.CreateCancelKeyboard())
}

// processTaskText в фсм состоянии ждет название таска
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

// processTaskDate в фсм состоянии ждет дату таски
func (h *TaskHandler) processTaskDate(c tele.Context, state *fsm.FSMData) error {
	if c.Callback() == nil {
		getDate := c.Text()
		_, err := keyboards.ParseDate(getDate)
		if err != nil {
			return c.Send(err.Error())
		}
		state.TaskDate = getDate
	} else {
		// получаю callback данные и убираю лишние символы
		rawDate := c.Callback().Data
		cleanDate := strings.Join(strings.Fields(rawDate), " ")
		state.TaskDate = cleanDate
	}
	fmt.Println("Мы тут", state)
	err := h.createTask(c, state)
	if err != nil {
		return c.Edit(err.Error())
	}

	return nil

}

// Complete закрытие задачи
func (h *TaskHandler) Complete(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	strTaskID := c.Callback().Data
	taskID, err := strconv.Atoi(strTaskID)
	if err != nil {
		return c.Edit(err.Error())
	}
	_, err = h.service.MarkDone(uint(taskID))
	if err != nil {
		return err
	}
	return c.Edit("✅ Задача завершена", keyboards.CreateMainKeyboard())
}

func (h *TaskHandler) createTask(c tele.Context, state *fsm.FSMData) error {
	taskDescription := state.TaskText
	taskDeadline, err := keyboards.ParseDate(state.TaskDate)
	if err != nil {
		return c.Send(err.Error())
	}

	task, err := h.service.Add(domain.Task{
		Description: taskDescription,
		Deadline:    taskDeadline,
	})
	if err != nil {
		return c.Send(err.Error())
	}
	err = h.fsmService.ClearState(h.ctx, c.Sender().ID)

	if err != nil {
		return c.Send(err.Error())
	}
	return c.Send(formatTask(task), keyboards.CreateMainKeyboard())
}
