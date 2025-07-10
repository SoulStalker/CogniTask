package handlers

import (
	"context"

	"github.com/SoulStalker/cognitask/internal/domain"
	"github.com/SoulStalker/cognitask/internal/keyboards"
	"github.com/SoulStalker/cognitask/internal/messages"
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
	fileType := c.Message().Media().MediaType()
	err := h.service.Create(domain.Media{Link: link, Type: fileType})
	if err != nil {
		return c.Send(err.Error(), keyboards.CreateMainKeyboard())
	}
	return c.Send(messages.BotMessages.FileSaved, keyboards.CreateMainKeyboard())
}

func (h *MediaHandler) Delete(c tele.Context) error {
	link := c.Message().Media().MediaFile().FileID
	err := h.service.Delete(domain.Media{Link: link})
	if err != nil {
		return c.Send(err.Error())
	}
	return c.Send(messages.BotMessages.FileDeleted, keyboards.CreateMainKeyboard())
}

func (h *MediaHandler) GetByLink(c tele.Context) error {
	return c.Send("Not implemented")
}

func (h *MediaHandler) Random(c tele.Context) error {
	media, err := h.service.Random()
	if err != nil {
		return c.Send(err.Error())
	}
	switch media.Type {
	case "photo":
		return c.Send(&tele.Photo{File: tele.File{FileID: media.Link}})
	case "video":
		return c.Send(&tele.Video{File: tele.File{FileID: media.Link}})
	default:
		return c.Send(messages.BotMessages.UnknownType)
	}
}
