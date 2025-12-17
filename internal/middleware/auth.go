package middleware

import (
	"cashly/internal/state"

	tb "gopkg.in/telebot.v3"
)

func Auth(password string) func(tb.HandlerFunc) tb.HandlerFunc {
	return func(next tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			uid := c.Sender().ID

			if state.IsAuthorized(uid) {
				return next(c)
			}

			if c.Text() != "" {
				if c.Text() == password {
					c.Delete()
					state.SetAuthorized(uid)
					return next(c)
				}
			}

			if c.Callback() != nil {
				_ = c.Respond(&tb.CallbackResponse{
					Text:      "🔐 Ти не авторизований!",
					ShowAlert: true,
				})
				return nil
			}

			return c.Send("🔐 Введи пароль для доступу:")
		}
	}
}
