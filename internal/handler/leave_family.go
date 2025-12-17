package handler

import (
	"cashly/internal/entity"
	"context"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) LeaveFamily(c tb.Context) error {
	inlineKeys := [][]tb.InlineButton{
		{BtnLeaveFamilyNo}, {BtnLeaveFamilyYes},
	}

	return c.Send("Ви дійсно хочете вийти з сім'ї?", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) ProcessLeaveFamily(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.LeaveFamily(ctx, f, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	h.eventCh <- entity.EventNotification{
		Type:        entity.EventLeavedFromFamily,
		RecipientID: f.CreatedBy,
		FamilyName:  f.Name,
		Data: map[string]any{
			"leaved_user_id": uid,
		},
	}

	h.bot.Edit(c.Message(), "Ти успішно вийшов з сім'ї.")

	return h.GoHome(c)
}

func (h *Handler) CancelLeaveFamily(c tb.Context) error {
	return c.Edit("Скасовано. Ти не вийшов з сім'ї.")
}
