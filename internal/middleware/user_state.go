package middleware

import (
	"cashly/internal/handler"
	"cashly/internal/state"

	tb "gopkg.in/telebot.v3"
)

func CheckUserState(goHome tb.HandlerFunc) func(tb.HandlerFunc) tb.HandlerFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			uid := c.Sender().ID

			us, ok := state.GetUserState(uid)
			if !ok || us.Family == nil {
				c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
				return goHome(c)
			}

			c.Set(handler.UfsKey, us.Family)
			return next(c)
		}
	}
}
