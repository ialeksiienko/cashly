package handler

import (
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"errors"
	"fmt"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateNewInviteCode(c tb.Context) error {
	userID := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	code, expiresAt, err := h.usecase.CreateNewInviteCode(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeNoPermission {
				return c.Send("У тебе немає прав на створення нового коду запрошення.")
			}
			if custErr.Code == errorsx.ErrCodeFailedToGenerateInviteCode {
				return c.Send("Не вдалося створити новий код запрошення. Спробуй пізніше.")
			}
		}
		return c.Send("Не вдалося створити код запрошення. Спробуй ще раз пізніше.")
	}

	return c.Send(fmt.Sprintf("Код запрошення: `%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
}
