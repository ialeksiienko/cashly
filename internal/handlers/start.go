package handlers

import (
	"cashly/internal/entity"
	"cashly/internal/state"
	"context"
	"log/slog"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) Start(c tb.Context) error {
	u := c.Sender()
	uid := u.ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := h.usecase.RegisterUser(ctx, &entity.User{
		ID:        uid,
		Username:  u.Username,
		Firstname: u.FirstName,
	})
	if err != nil {
		h.logger.Error("failed to save user", slog.Int("user_id", int(uid)), slog.String("err", err.Error()))
		return c.Send("Сталася помилка при зберіганні данних користувача. Спробуй пізніше.")
	}

	if _, exists := state.GetUserState(uid); !exists {
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
