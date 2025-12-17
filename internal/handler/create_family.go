package handler

import (
	"cashly/internal/state"
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи назву нової сім'ї:")

	state.SetTextState(c.Sender().ID, state.WaitingFamilyName)

	return nil
}

func (h *Handler) processFamilyCreation(c tb.Context, name string) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if utf8.RuneCountInString(name) > 20 {
		return c.Send("Назва сім'ї не має містити більше 20 символів.")
	}

	_, code, expiresAt, err := h.usecase.CreateFamily(ctx, name, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Сім'я `%s` створена. Код запрошення:\n\n`%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", name, code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}
