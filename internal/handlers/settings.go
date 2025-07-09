package handlers

import (
	"fmt"

	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

type SettingsHandler struct {
	service usecase.SettingsService
}

func NewSettingsHandler(uc usecase.SettingsService) *SettingsHandler {
	return &SettingsHandler{
		service: uc,
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
