package handler

import (
	"cashly/internal/entity"
	"cashly/internal/state"
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) JoinFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи код запрошення:")

	state.SetTextState(c.Sender().ID, state.WaitingFamilyCode)

	return nil
}

func (h *Handler) processFamilyJoin(c tb.Context, code string) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if utf8.RuneCountInString(code) != 6 {
		return c.Send("Код запрошення має містити 6 символів.")
	}

	f, err := h.usecase.JoinFamily(ctx, code, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	h.eventCh <- entity.EventNotification{
		Type:        entity.EventJoinedFamily,
		FamilyName:  f.Name,
		RecipientID: f.CreatedBy,
		Data: map[string]any{
			"joined_user_id": uid,
		},
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Успішно приєднався до сім'ї! Назва - %s", f.Name), &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
