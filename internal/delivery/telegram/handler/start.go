package handler

import (
	"cashly/internal/entity"
	"cashly/internal/session"
	"context"
	"log/slog"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) Start(c tb.Context) error {
	user := c.Sender()
	userID := user.ID
	ctx := context.Background()

	_, err := h.usecase.RegisterUser(ctx, &entity.User{
		ID:        userID,
		Username:  user.Username,
		Firstname: user.FirstName,
	})
	if err != nil {
		h.sl.Error("failed to save user", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
		return c.Send("Сталася помилка при зберіганні данних користувача. Спробуй пізніше.")
	}

	if _, exists := session.GetUserState(userID); !exists {
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

	return c.Send("Привіт! Цей бот допоможе дізнатися рахунок на карті Monobank.\n\nВибери один з варіантів на клавіатурі.", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
