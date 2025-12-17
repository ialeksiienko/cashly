package middleware

import (
	"encoding/json"
	"log/slog"
	"os"

	tb "gopkg.in/telebot.v3"
)

type Family struct {
	Firstname string `json:"firstname"`
	ID        int64  `json:"id"`
}

var allowedUsers map[int64]struct{}

func init() {
	content, err := os.ReadFile("family.json")
	if err != nil {
		slog.Error("failed to read family.json", slog.String("err", err.Error()))
		os.Exit(1)
	}

	var fam []Family

	if unmarshalErr := json.Unmarshal(content, &fam); unmarshalErr != nil {
		slog.Error("failed to unmarshal family.json", slog.String("err", unmarshalErr.Error()))
		os.Exit(1)
	}

	allowedUsers = make(map[int64]struct{})
	for _, f := range fam {
		allowedUsers[f.ID] = struct{}{}
	}
}

func isAllowed(id int64) bool {
	_, ok := allowedUsers[id]
	return ok
}

func CheckAllowedUsers(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		userID := c.Sender().ID

		if !isAllowed(userID) {
			return c.Send("У тебе немає прав користуватися ботом, зв'яжись з адміністратором.")
		}

		return next(c)
	}
}
