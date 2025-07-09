package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

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
}

func NewSettingsHandler(fsmService *fsm.FSMService, uc usecase.SettingsService, ctx context.Context) *SettingsHandler {
	return &SettingsHandler{
		fsmService: fsmService,
		service:    uc,
		ctx:        ctx,
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
	var sets string

	// пока так замокаю

	sets += fmt.Sprintf("DeleteAfterDays: %d\n", settings.DeleteAfterDays)
	sets += fmt.Sprintf("NotificationHours: %d\n", settings.NotificationHours)
	sets += fmt.Sprintf("NotifyFrom: %d\n", settings.NotifyFrom)
	sets += fmt.Sprintf("NotifyTo: %d\n", settings.NotifyTo)
	sets += fmt.Sprintf("RandomHour: %d\n", settings.RandomHour)

	return c.Edit(sets, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) SetDeleteDays(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}

	userID := c.Sender().ID
	if err := h.fsmService.ClearState(h.ctx, userID); err != nil {
		log.Printf("Failed to clear state: %v", err)
	}

	state := &fsm.FSMData{
		State: fsm.StateDeleteAfterDays,
	}

	if err := h.fsmService.SetState(h.ctx, userID, state); err != nil {
		log.Printf("Failed to set state: %v", err)
		c.Edit(messages.BotMessages.ErrorSomeError)
	}
	return c.Edit("Выбери число: ", keyboards.CreateHoursKeyboard(4))
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
