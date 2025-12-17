package handlers

import (
	"cashly/internal/state"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) GoHome(c tb.Context) error {
	uid := c.Sender().ID

	state.DeleteUserState(uid)

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
