package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"
	"github.com/SoulStalker/cognitask/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

type SettingsHandler struct {
	fsmService *fsm.FSMService
	service    usecase.SettingsService
	ctx        context.Context
	ch         chan time.Duration
}

func NewSettingsHandler(fsmService *fsm.FSMService, uc usecase.SettingsService, ctx context.Context, ch chan time.Duration) *SettingsHandler {
	return &SettingsHandler{
		fsmService: fsmService,
		service:    uc,
		ctx:        ctx,
		ch:         ch,
	}
}

func (h *SettingsHandler) CanHandle(state string) bool {
	return state == fsm.StateDeleteAfterDays ||
		state == fsm.StateNotificationHours ||
		state == fsm.StateNotifyFrom ||
		state == fsm.StateNotifyTo ||
		state == fsm.StateRandom
}

func (h *SettingsHandler) Handle(c tele.Context, data *fsm.FSMData) error {
	switch data.State {
	case fsm.StateDeleteAfterDays:
		return h.processDeleteDays(c)
	case fsm.StateNotificationHours:
		return h.processNotificationHours(c)
	case fsm.StateNotifyFrom:
		return h.processNotifyFrom(c)
	case fsm.StateNotifyTo:
		return h.processNotifyTo(c)
	case fsm.StateRandom:
		return h.processRandomHour(c)
	default:
		return c.Send("unknown callback")
	}
}

func (h *SettingsHandler) Settings(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}
	settings, err := h.service.All()
	if err != nil {
		return c.Edit(err.Error())
	}
	currentSettings := "⚙️ Текущие настройки:\n--------------------------\n\n"

	currentSettings += fmt.Sprintf("🗑️ Авто-удаление выполненных задач через дней: %d\n\n", settings.DeleteAfterDays)
	currentSettings += fmt.Sprintf("⏰ Период уведомлений (часов): %d\n\n", settings.NotificationHours)
	currentSettings += fmt.Sprintf("📅 Начало уведомлений в: %d\n\n", settings.NotifyFrom)
	currentSettings += fmt.Sprintf("📅 Конец уведомлений в: %d\n\n", settings.NotifyTo)
	currentSettings += fmt.Sprintf("💡 Мотиватор в: %d\n\n", settings.RandomHour)

	currentSettings += "Можешь изменить настройки по кнопкам ниже:"

	return c.Edit(currentSettings, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) SetDeleteDays(c tele.Context) error {
	err := h.setState(c, fsm.StateDeleteAfterDays)
	if err != nil {
		c.Edit(err.Error())
	}
	return nil
}

func (h *SettingsHandler) SetNotificationHours(c tele.Context) error {
	err := h.setState(c, fsm.StateNotificationHours)
	if err != nil {
		c.Edit(err.Error())
	}
	return nil
}

func (h *SettingsHandler) SetNotifyFrom(c tele.Context) error {
	err := h.setState(c, fsm.StateNotifyFrom)
	if err != nil {
		c.Edit(err.Error())
	}
	return nil
}

func (h *SettingsHandler) SetNotifyTo(c tele.Context) error {
	err := h.setState(c, fsm.StateNotifyTo)
	if err != nil {
		c.Edit(err.Error())
	}
	return nil
}

func (h *SettingsHandler) SetRandomHour(c tele.Context) error {
	err := h.setState(c, fsm.StateRandom)
	if err != nil {
		c.Edit(err.Error())
	}
	return nil
}

func (h *SettingsHandler) processDeleteDays(c tele.Context) error {
	rawDays := c.Callback().Data
	cleanDays := strings.Join(strings.Fields(rawDays), " ")
	deleteDays, err := strconv.Atoi(cleanDays)
	if err != nil {
		return c.Edit(err.Error())
	}

	err = h.service.SetDeleteDays(uint(deleteDays))
	if err != nil {
		c.Edit(err.Error())
	}

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	return c.Edit(cleanDays, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) setState(c tele.Context, newState string) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	userID := c.Sender().ID
	if err := h.fsmService.ClearState(h.ctx, userID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	state := &fsm.FSMData{
		State: newState,
	}

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to set state: %v", err)
		c.Edit(messages.BotMessages.ErrorSomeError)
	}
	return c.Edit("Выбери число: ", keyboards.CreateHoursKeyboard(4))
}

func (h *SettingsHandler) processNotificationHours(c tele.Context) error {
	rawHours := c.Callback().Data
	cleanHours := strings.Join(strings.Fields(rawHours), " ")
	hours, err := strconv.Atoi(cleanHours)
	if err != nil {
		return c.Edit(err.Error())
	}

	err = h.service.SetNotificationHours(uint(hours))
	if err != nil {
		c.Edit(err.Error())
	}

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}
	h.ch <- time.Duration(hours) * time.Second // минуты для тестов
	return c.Edit(cleanHours, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) processNotifyFrom(c tele.Context) error {
	rawHours := c.Callback().Data
	cleanHours := strings.Join(strings.Fields(rawHours), " ")
	hours, err := strconv.Atoi(cleanHours)
	if err != nil {
		return c.Edit(err.Error())
	}

	err = h.service.SetNotifyFrom(uint(hours))
	if err != nil {
		c.Edit(err.Error())
	}

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	return c.Edit(cleanHours, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) processNotifyTo(c tele.Context) error {
	rawHours := c.Callback().Data
	cleanHours := strings.Join(strings.Fields(rawHours), " ")
	hours, err := strconv.Atoi(cleanHours)
	if err != nil {
		return c.Edit(err.Error())
	}

	err = h.service.SetNotifyTo(uint(hours))
	if err != nil {
		c.Edit(err.Error())
	}

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	return c.Edit(cleanHours, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) processRandomHour(c tele.Context) error {
	rawHour := c.Callback().Data
	cleanHour := strings.Join(strings.Fields(rawHour), " ")
	hours, err := strconv.Atoi(cleanHour)
	if err != nil {
		return c.Edit(err.Error())
	}

	err = h.service.SetRandomHour(uint(hours))
	if err != nil {
		c.Edit(err.Error())
	}

	if err := h.fsmService.ClearState(h.ctx, c.Sender().ID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	return c.Edit(cleanHour, keyboards.CreateSettingsKeyboard())
}
