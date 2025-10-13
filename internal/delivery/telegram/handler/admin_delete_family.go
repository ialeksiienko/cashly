package handler

import (
	"context"
	"errors"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) DeleteFamily(c tb.Context) error {
	inlineKeys := [][]tb.InlineButton{
		{BtnFamilyDeleteNo}, {BtnFamilyDeleteYes},
	}

	return c.Send("Дійсно хочеш видалити сім'ю?", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) ProcessFamilyDeletion(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	h.bot.Delete(c.Message())

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Edit(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.DeleteFamily(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeNoPermission {
				return c.Edit("У тебе немає прав на видалення.")
			}
		}
		return c.Edit("Не вдалося видалити сім'ю. Спробуйте ще раз пізніше.")
	}

	h.bot.Edit(c.Message(), "Сім'ю успішно видалено.")

	return h.GoHome(c)
}

func (h *Handler) CancelFamilyDeletion(c tb.Context) error {
	return c.Edit("Скасовано. Сім’ю не було видалено.")
}
