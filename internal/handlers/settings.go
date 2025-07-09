package handlers

import (
	"context"
	"fmt"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/keyboards"
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
		return h.processDeleteDays(c, data)
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

	sets += fmt.Sprintf("%d\n", settings.DeleteAfterDays)
	sets += fmt.Sprintf("%d\n", settings.NotificationHours)
	sets += fmt.Sprintf("%d\n", settings.NotifyFrom)
	sets += fmt.Sprintf("%d\n", settings.NotifyTo)
	sets += fmt.Sprintf("%d\n", settings.RandomHour)

	return c.Edit(sets, keyboards.CreateSettingsKeyboard())
}

func (h *SettingsHandler) SetDeleteDays(c tele.Context) error {
	err := c.Respond()
	if err != nil {
		return c.Edit(err.Error())
	}
	// err = h.service.SetDeleteDays(100)
	// if err != nil {
	// 	return c.Edit(err.Error())
	// }
	// return c.Edit("ok", keyboards.CreateSettingsKeyboard())

	return c.Edit("Выбери нужный час: ", keyboards.CreateHoursKeyboard(4))
}

func (h *SettingsHandler) processDeleteDays(c tele.Context, data *fsm.FSMData) error {
	return nil
}
