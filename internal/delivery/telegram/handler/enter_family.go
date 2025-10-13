package handler

import (
	"context"
	"errors"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) EnterMyFamily(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	families, err := h.usecase.GetFamiliesByUserID(ctx, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeUserHasNoFamily {
				inlineKeys := [][]tb.InlineButton{
					{BtnCreateFamily}, {BtnJoinFamily},
				}

				return c.Edit("Привіт! У вас поки немає жодної сім'ї. Створи або приєднайся.", &tb.ReplyMarkup{
					InlineKeyboard: inlineKeys,
				})
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	session.SetUserPageState(userID, &session.UserPageState{
		Page:     0,
		Families: families,
	})

	return showFamilyListPage(c, families, 0)
}
