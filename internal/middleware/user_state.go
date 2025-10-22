package middleware

import (
	"cashly/internal/session"

	tb "gopkg.in/telebot.v3"
)

func CheckUserState(goHome tb.HandlerFunc) func(tb.HandlerFunc) tb.HandlerFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			userID := c.Sender().ID

			userState, exists := session.GetUserState(userID)
			if !exists || userState.Family == nil {
				c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
				return goHome(c)
			}

			c.Set("user_state", userState)
			return next(c)
		}
	}
}
