package handler

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) JoinFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи код запрошення:")

	session.SetTextState(c.Sender().ID, session.StateWaitingFamilyCode)

	return nil
}

func (h *Handler) processFamilyJoin(c tb.Context, code string) error {
	userID := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if utf8.RuneCountInString(code) != 6 {
		return c.Send("Код запрошення має містити 6 символів.")
	}

	family, err := h.usecase.JoinFamily(ctx, code, userID)
	if err != nil {
		switch e := err.(type) {
		case *errorsx.CustomError[time.Time]:
			if e.Code == errorsx.ErrCodeFamilyCodeExpired {
				return c.Send(fmt.Sprintf("Код запрошення не дійсний, закінчився - %s (час за Гринвічем, GMT)", e.Data.Format("02.01.2006 о 15:04")))
			}
		case *errorsx.CustomError[struct{}]:
			if e.Code == errorsx.ErrCodeFamilyNotFound {
				return c.Send("Сім'ю з цим кодом запрошення не знайдено.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	h.eventCh <- &entity.EventNotification{
		Event:       entity.EventJoinedFamily,
		FamilyName:  family.Name,
		RecipientID: family.CreatedBy,
		Data: map[string]any{
			"joined_user_id": userID,
		},
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Успішно приєднався до сім'ї! Назва - %s", family.Name), &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
