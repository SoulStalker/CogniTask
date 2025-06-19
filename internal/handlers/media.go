package handlers

import (
	"context"
	"fmt"

	"github.com/SoulStalker/cognitask/internal/usecase"
	tele "gopkg.in/telebot.v3"
)

type MediaHandler struct {
	service *usecase.MediaService
	ctx     context.Context
}

func NewMediaHandler(service *usecase.MediaService, ctx context.Context) *MediaHandler {
	return &MediaHandler{
		service: service,
		ctx:     ctx,
	}
}

func (h *MediaHandler) Create(c tele.Context) error {
	link := c.Message().Media().MediaFile().FileID
	fmt.Println(link)
	return c.Send(link)
}

func (h *MediaHandler) Delete(c tele.Context) error {
	return c.Send("Not implemented")
}

func (h *MediaHandler) GetByLink(c tele.Context) error {
	return c.Send("Not implemented")
}

func (h *MediaHandler) Random(c tele.Context) error {
	return c.Send("Not implemented")
}
