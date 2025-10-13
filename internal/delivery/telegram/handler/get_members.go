package handler

import (
	"context"
	"errors"
	"fmt"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) GetMembers(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	members, err := h.usecase.GetFamilyMembersInfo(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeFamilyHasNoMembers {
				return c.Send("У твоїй сім'ї поки немає учасників.")
			}
		}
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

	c.Send("📋 Список учасників сім'ї:\n")

	for _, member := range members {
		role := "Учасник"
		if member.IsAdmin {
			role = "Адміністратор"
		}

		userLabel := ""
		if member.IsCurrent {
			userLabel = " (це ви)"
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s\n- Роль: %s\n- ID: %d",
			member.Firstname,
			member.Username,
			userLabel,
			role,
			member.ID,
		)

		isAdmin := userID == us.Family.CreatedBy

		if !member.IsCurrent && isAdmin {
			btn := tb.InlineButton{
				Unique: "delete_member",
				Text:   "🗑 Видалити",
				Data:   strconv.FormatInt(member.ID, 10),
			}

			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{
				{btn},
			}

			c.Send(text, markup)
		} else {
			c.Send(text)
		}
	}

	return c.Send(fmt.Sprintf("Всього учасників: %d", len(members)))
}
