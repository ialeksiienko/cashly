package handler

import (
	"cashly/internal/entity"
	"context"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) DeleteFamily(c tb.Context) error {
	ib := [][]tb.InlineButton{
		{BtnFamilyDeleteNo}, {BtnFamilyDeleteYes},
	}

	return c.Send("Дійсно хочеш видалити сім'ю?", &tb.ReplyMarkup{
		InlineKeyboard: ib,
	})
}

func (h *Handler) ProcessFamilyDeletion(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	h.bot.Delete(c.Message())

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.DeleteFamily(ctx, f, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	h.bot.Edit(c.Message(), "Сім'ю успішно видалено.")

	return h.GoHome(c)
}

func (h *Handler) CancelFamilyDeletion(c tb.Context) error {
	return c.Edit("Скасовано. Сім’ю не було видалено.")
}
