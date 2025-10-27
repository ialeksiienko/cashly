package handler

import (
	"cashly/internal/entity"
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"fmt"
	"strconv"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) DeleteMember(c tb.Context) error {
	data := c.Callback().Data

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	member, err := h.usecase.GetUserByID(ctx, memberID)
	if err != nil {
		return c.Send(ErrInternalServerForUser.Error())
	}

	inlineKeys := [][]tb.InlineButton{
		{tb.InlineButton{Unique: "delete_member_no", Text: "❌ Ні", Data: strconv.FormatInt(member.ID, 10)}}, {tb.InlineButton{Unique: "delete_member_yes", Text: "✅ Так", Data: strconv.FormatInt(member.ID, 10)}},
	}

	return c.Edit(
		fmt.Sprintf("Дійсно хочеш видалити учасника `%s`?", member.Firstname),
		&tb.SendOptions{
			ParseMode:   tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{InlineKeyboard: inlineKeys},
		},
	)
}

func (h *Handler) ProcessMemberDeletion(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Edit("Некоректний ID користувача.")
	}

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Edit(ErrUnableToGetUserState.Error())
	}

	removeErr := h.usecase.RemoveMember(ctx, us.Family.ID, userID, memberID)
	if removeErr != nil {
		switch e := err.(type) {
		case *errorsx.CustomError[struct{}]:
			if e.Code == errorsx.ErrCodeNoPermission {
				return c.Edit("У тебе немає прав на видалення.")
			}
			if e.Code == errorsx.ErrCodeCannotRemoveSelf {
				return c.Edit("Ти не можеш видалити себе.")
			}
		}
		return c.Edit("Не вдалося видалити користувача з сім'ї. Спробуй ще раз пізніше.")
	}

	h.eventCh <- &entity.EventNotification{
		Event:       entity.EventDeletedFromFamily,
		RecipientID: memberID,
		FamilyName:  us.Family.Name,
	}

	return c.Delete()
}

func (h *Handler) CancelMemberDeletion(c tb.Context) error {
	data := c.Callback().Data

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Edit("Некоректний ID користувача.")
	}

	return c.Edit("Скасовано. Учасника не було видалено.", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{{Unique: "go_back_delete_member", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(memberID), 10)}}}})
}
