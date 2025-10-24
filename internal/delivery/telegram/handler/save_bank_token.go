package handler

import (
	"cashly/internal/session"
	"cashly/internal/validate"
	"context"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) SaveUserBankToken(c tb.Context) error {
	button := tb.InlineButton{
		Unique: "mono_link",
		Text:   "Посилання",
		URL:    "https://api.monobank.ua/",
	}

	inlineKeys := [][]tb.InlineButton{
		{button},
	}

	h.bot.Send(c.Sender(), "Перейди по посиланню знизу та відправ свій токен в цей чат.", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})

	session.SetTextState(c.Sender().ID, session.StateWaitingBankToken)

	return nil
}

func (h *Handler) processUserBankToken(c tb.Context) error {
	userID := c.Sender().ID
	token := c.Message().Text
	ctx := context.Background()

	us, exists := session.GetUserState(userID)
	if !exists || us.Family == nil {
		c.Send("Ти не увійшов в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	valid := validate.IsValidBankToken(token)
	if !valid {
		return c.Send("Неправильний формат токена.")
	}

	_, saveErr := h.usecase.SaveBankToken(ctx, us.Family.ID, userID, token)
	if saveErr != nil {
		return c.Send("Не вдалося зберегти токен. Спробуй пізніше.")
	}

	isAdmin := us.Family.CreatedBy == userID

	rows := generateFamilyMenu(isAdmin, true)

	menu.Reply(rows...)

	c.Delete()

	return c.Send("Ти успішно зберіг токен для цієї сім'ї.", menu)
}
