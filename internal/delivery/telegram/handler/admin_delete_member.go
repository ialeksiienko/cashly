package handler

import (
	"context"
	"fmt"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) DeleteMember(c tb.Context) error {
	data := c.Callback().Data
	ctx := context.Background()

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	member, err := h.usecase.GetUserByID(ctx, memberID)
	if err != nil {
		return c.Send(ErrInternalServerForUser.Error())
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnMemberDeleteNo}, {tb.InlineButton{Unique: "delete_member_yes", Text: "✅ Так", Data: strconv.FormatInt(member.ID, 10)}},
	}

	return c.Send(
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
	ctx := context.Background()

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

	h.bot.Edit(c.Message(), "Учасника успішно видалено. Оновлюю список...")

	h.bot.Send(c.Sender(), "── 🔹 Оновлення списку 🔹 ──")

	return h.GetMembers(c)
}

func (h *Handler) CancelMemberDeletion(c tb.Context) error {
	return c.Edit("Скасовано. Учасника не було видалено.")
}
