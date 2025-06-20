package handlers

import (
	"context"

	"github.com/SoulStalker/cognitask/internal/domain"
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
	err := h.service.Create(domain.Media{Link: link})
	if err != nil {
		c.Send(err.Error())
	}
	return c.Send(link)
}

func (h *MediaHandler) Delete(c tele.Context) error {
	link := c.Message().Media().MediaFile().FileID
	err := h.service.Delete(domain.Media{Link: link})
	if err != nil {
		return c.Send(err.Error())
	}
	return c.Send("Deleted")
}

func (h *MediaHandler) GetByLink(c tele.Context) error {
	return c.Send("Not implemented")
}

func (h *MediaHandler) Random(c tele.Context) error {
	media, err := h.service.Random()
	if err != nil {
		return c.Send(err.Error())
	}
	photo := &tele.Photo{File: tele.File{FileID: media.Link}}
	return c.Send(photo)
}
