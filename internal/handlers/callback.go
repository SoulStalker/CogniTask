package handlers

import (
	"context"
	"fmt"
	"log"

	"github.com/SoulStalker/cognitask/internal/fsm"
	"github.com/SoulStalker/cognitask/internal/messages"
	tele "gopkg.in/telebot.v3"
)

type CallbackHandler interface {
	CanHandle(state string) bool
	Handle(c tele.Context, data *fsm.FSMData) error
}

type CallbackRouter struct {
	handlers []CallbackHandler
	fsmSvc   *fsm.FSMService
	ctx      context.Context
}

func NewCallbackRouter(handlers []CallbackHandler, fsmSvc *fsm.FSMService, ctx context.Context) *CallbackRouter {
	return &CallbackRouter{
		handlers: handlers,
		fsmSvc: fsmSvc,
		ctx: ctx,
	}
}

func (r *CallbackRouter) Handle(c tele.Context) error {
	userID := c.Sender().ID

	data, err := r.fsmSvc.GetState(r.ctx, userID)
	if err != nil {
		log.Printf("Failed to get state: %v", err)
		return c.Send(messages.BotMessages.ErrorTryAgain)
	}

	for _, h := range r.handlers {
		if h.CanHandle(data.State) {
			return h.Handle(c, data)
		}
	}

	return c.Send(fmt.Sprintf("Wrong callback int state %s", data))
}
