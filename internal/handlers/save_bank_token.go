package handlers

import (
	"cashly/internal/state"
	"cashly/internal/validate"
	"context"
	"time"

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

	state.SetTextState(c.Sender().ID, state.WaitingBankToken)

	return nil
}

func (h *Handler) processUserBankToken(c tb.Context) error {
	uid := c.Sender().ID
	token := c.Message().Text

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	us, exists := state.GetUserState(uid)
	if !exists || us.Family == nil {
		c.Send("Ти не увійшов в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	valid := validate.BankToken(token)
	if !valid {
		return c.Edit("Неправильний формат токена.")
	}

	_, serr := h.usecase.SaveBankToken(ctx, us.Family.ID, uid, token)
	if serr != nil {
		return c.Edit("Не вдалося зберегти токен. Спробуй пізніше.")
	}

	isAdmin := us.Family.CreatedBy == uid

	rows := generateFamilyMenu(isAdmin, true)

	menu.Reply(rows...)

	c.Delete()

	return c.Send("Ти успішно зберіг токен для цієї сім'ї.", menu)
}
