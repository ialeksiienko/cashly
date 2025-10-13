package handler

import (
	"context"
	"monofamily/internal/session"

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
	userID := c.Sender().ID
	ctx := context.Background()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Edit(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.DeleteUserBankToken(ctx, us.Family.ID, userID)
	if err != nil {
		return c.Edit(ErrInternalServerForUser.Error())
	}

	rows := generateFamilyMenu(us.Family.CreatedBy == userID, false)

	menu.Reply(rows...)

	return c.Send("Токен успішно видалено.", menu)
}

func (h *Handler) CancelRemoveBankToken(c tb.Context) error {
	return c.Edit("Скасовано. Токен не було видалено.")
}
