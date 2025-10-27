package handler

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"errors"
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
	userID := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.LeaveFamily(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeCannotRemoveSelf {
				return c.Edit("Адміністратор не може вийти з сім'ї.")
			}
		}
		return c.Edit("Не вдалося вийти з сім'ї. Спробуйте ще раз пізніше.")
	}

	h.eventCh <- &entity.EventNotification{
		Event:       entity.EventLeavedFromFamily,
		RecipientID: us.Family.CreatedBy,
		FamilyName:  us.Family.Name,
		Data: map[string]any{
			"leaved_user_id": userID,
		},
	}

	h.bot.Edit(c.Message(), "Ти успішно вийшов з сім'ї.")

	return h.GoHome(c)
}

func (h *Handler) CancelLeaveFamily(c tb.Context) error {
	return c.Edit("Скасовано. Ти не вийшов з сім'ї.")
}
