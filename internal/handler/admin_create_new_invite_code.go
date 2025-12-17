package handler

import (
	"cashly/internal/entity"
	"context"
	"fmt"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateNewInviteCode(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	code, expiresAt, err := h.usecase.CreateNewInviteCode(ctx, f, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	return c.Send(fmt.Sprintf("Код запрошення: `%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
}
