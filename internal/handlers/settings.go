package handlers

import (
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
	return c.Edit("Settings", keyboards.CreateSettingsKeyboard())
}
