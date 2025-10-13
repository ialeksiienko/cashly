package handler

import (
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) GoHome(c tb.Context) error {
	userID := c.Sender().ID

	session.DeleteUserState(userID)

	{
		msg, _ := h.bot.Send(c.Sender(), ".", &tb.SendOptions{
			ReplyMarkup: &tb.ReplyMarkup{
				RemoveKeyboard: true,
			},
		})

		h.bot.Delete(msg)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnCreateFamily}, {BtnJoinFamily}, {BtnEnterMyFamily},
	}

	return c.Send("Вибери один з варіантів на клавіатурі.", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
