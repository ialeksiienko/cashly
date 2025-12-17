package handler

import (
	"cashly/internal/entity"
	"context"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) RemoveBankToken(c tb.Context) error {
	inlineKeys := [][]tb.InlineButton{
		{BtnRemoveBankTokenNo}, {BtnRemoveBankTokenYes},
	}

	return c.Send("Впевнений що хочеш видалити токен монобанку?", &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{InlineKeyboard: inlineKeys}},
	)
}

func (h *Handler) ProcessRemoveBankToken(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.DeleteUserBankToken(ctx, f.ID, uid)
	if err != nil {
		return c.Edit(ErrInternalServerForUser.Error())
	}

	rows := generateFamilyMenu(f.CreatedBy == uid, false)

	menu.Reply(rows...)

	return c.Send("Токен успішно видалено.", menu)
}

func (h *Handler) CancelRemoveBankToken(c tb.Context) error {
	return c.Edit("Скасовано. Токен не було видалено.")
}
