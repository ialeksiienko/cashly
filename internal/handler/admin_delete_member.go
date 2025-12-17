package handler

import (
	"cashly/internal/entity"
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
	uid := c.Sender().ID
	d := c.Callback().Data

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mid, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return c.Edit("Некоректний ID користувача.")
	}

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	rerr := h.usecase.RemoveMember(ctx, f.ID, uid, mid)
	if rerr != nil {
		return c.Send(mapErrorToMessage(rerr))
	}

	h.eventCh <- entity.EventNotification{
		Type:        entity.EventDeletedFromFamily,
		RecipientID: mid,
		FamilyName:  f.Name,
	}

	return c.Delete()
}

func (h *Handler) CancelMemberDeletion(c tb.Context) error {
	d := c.Callback().Data

	mid, err := strconv.ParseInt(d, 10, 64)
	if err != nil {
		return c.Edit("Некоректний ID користувача.")
	}

	return c.Edit("Скасовано. Учасника не було видалено.", &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{{Unique: "go_back_delete_member", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(mid), 10)}}}})
}
