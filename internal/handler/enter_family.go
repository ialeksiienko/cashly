package handler

import (
	"cashly/internal/state"
	"context"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) EnterMyFamily(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, err := h.usecase.GetFamiliesByUserID(ctx, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	state.SetUserPageState(uid, &state.UserPage{
		Page:     0,
		Families: f,
	})

	return showFamilyListPage(c, f, 0)
}
