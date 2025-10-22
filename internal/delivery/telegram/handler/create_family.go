package handler

import (
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"errors"
	"fmt"
	"unicode/utf8"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи назву нової сім'ї:")

	session.SetTextState(c.Sender().ID, session.StateWaitingFamilyName)

	return nil
}

func (h *Handler) processFamilyCreation(c tb.Context, familyName string) error {
	userID := c.Sender().ID
	ctx := context.Background()

	if utf8.RuneCountInString(familyName) > 20 {
		return c.Send("Назва сім'ї не має містити більше 20 символів.")
	}

	_, code, expiresAt, err := h.usecase.CreateFamily(ctx, familyName, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeFailedToGenerateInviteCode {
				return c.Send("Не вдалося створити новий код запрошення. Спробуй пізніше.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Сім'я `%s` створена. Код запрошення:\n\n`%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", familyName, code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
