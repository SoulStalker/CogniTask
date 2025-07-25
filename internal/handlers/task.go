package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/mappers"
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

func (h *TaskHandler) CanHandle(state string) bool {
	return state == fsm.StateWaitingTaskText ||
		state == fsm.StateWaitingTaskDate ||
		state == fsm.StateWaitingTaskCategory
}

func (h *TaskHandler) Handle(c tele.Context, data *fsm.FSMData) error {
	switch data.State {
	case fsm.StateWaitingTaskText:
		return h.processTaskText(c, data)
	case fsm.StateWaitingTaskDate:
		return h.processTaskDate(c, data)
	default:
		return c.Send("unknown callback")
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
		return c.Edit(messages.BotMessages.NoOpenTasks, keyboards.CreateMainKeyboard())
	}

	rows := mappers.FormatTaskList(tasks)

	return c.Edit(messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})
}

// All выводит список всех задач
func (h *TaskHandler) All(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return err
	}

	tasks, err := h.service.All()
	if err != nil {
		return c.Send(err.Error())
	}
	if len(tasks) == 0 {
		return c.Edit(messages.BotMessages.NoTasks, keyboards.CreateMainKeyboard())
	}

	rows := mappers.FormatTaskList(tasks)

	return c.Edit(messages.BotMessages.YourTasks, &tele.ReplyMarkup{InlineKeyboard: rows})
}

// Add  хендлер для обработки команды  Add
func (h *TaskHandler) Add(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	userID := c.Sender().ID

	if err := h.fsmService.ClearState(h.ctx, userID); err != nil {
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

// textToTask если бот получил просто текст, можем из него сделать задачу
func (h *TaskHandler) textToTask(c tele.Context, state *fsm.FSMData) error {
	userID := c.Sender().ID
	state.TaskText = c.Text()
	state.State = fsm.StateWaitingTaskDate

	fmt.Println(state.State)

	// Надо сохранить новое состояние
	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to update state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}
	return c.Send(messages.BotMessages.InputNewDate, keyboards.GetDateSelectionKeyboard())

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
	return c.Edit(messages.BotMessages.TaskCompleted, keyboards.CreateMainKeyboard())
}

// Delete удаление задачи
func (h *TaskHandler) Delete(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	strTaskID := c.Callback().Data

	taskID, err := strconv.Atoi(strTaskID)
	if err != nil {
		return c.Edit(err.Error())
	}
	err = h.service.Delete(uint(taskID))
	if err != nil {
		return err
	}

	return c.Edit(messages.BotMessages.TaskDeleted, keyboards.CreateMainKeyboard())
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
	return c.Send(mappers.FormatTask(task), keyboards.CreateMainKeyboard())
}

func (h *TaskHandler) SelectDate(c tele.Context) error {
	// Создаем карту (map) для перевода месяцев
	monthTranslations := map[time.Month]string{
		time.January:   "Январь",
		time.February:  "Февраль",
		time.March:     "Март",
		time.April:     "Апрель",
		time.May:       "Май",
		time.June:      "Июнь",
		time.July:      "Июль",
		time.August:    "Август",
		time.September: "Сентябрь",
		time.October:   "Октябрь",
		time.November:  "Ноябрь",
		time.December:  "Декабрь",
	}

	year := time.Now().Year()
	month := monthTranslations[time.Now().Month()]
	markup := keyboards.BuildKeyboard(year, time.Now().Month())
	title := fmt.Sprintf("%s %d", month, year)
	return c.Send(title, markup)
}

// func sendCalendar(b *tele.Bot, c tele.Context, year int, month time.Month) error {
// 	markup := keyboards.BuildKeyboard(year, month)
// 	title := fmt.Sprintf("%s %d", month.String(), year)
// 	return c.Send(c.Sender(), title, markup)
// }

// func editCalendar(b *tele.Bot, c tele.Context, year int, month time.Month) error {
// 	markup := keyboards.BuildKeyboard(year, month)
// 	title := fmt.Sprintf("%s %d", month.String(), year)
// 	return c.EditCaption("c.Message()", title, markup)
// }

// Функция преобразует строку формата "day|DAY|YYYY|MM|DD" в "DD.MM.YYYY"
func mapDateString(input string) (string, error) {
	parts := strings.Split(input, "|")
	if len(parts) != 5 {
		return "", fmt.Errorf("неверный формат строки, ожидается day|DAY|YYYY|MM|DD")
	}

	year := parts[2]  // YYYY
	month := parts[3] // MM
	day := parts[4]   // DD

	// Проверяем, что компоненты даты корректны (можно добавить дополнительные проверки)
	if len(year) != 4 || len(month) != 2 || len(day) != 2 {
		return "", fmt.Errorf("некорректный формат даты")
	}

	return fmt.Sprintf("%s.%s.%s", day, month, year), nil
}