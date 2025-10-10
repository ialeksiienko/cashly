package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) ViewBalance(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	members, err := h.usecase.GetFamilyMembersInfo(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeFamilyHasNoMembers {
				return c.Send("У вашій сім'ї поки немає учасників.")
			}
		}
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

	c.Send("📋 Вибери учасника для перевірки балансу:\n")

	for _, member := range members {
		role := "Учасник"
		if member.IsAdmin {
			role = "Адміністратор"
		}

		btn := tb.InlineButton{}
		status := ""
		if !member.HasToken {
			status = " (користувач не додав токен)"
		} else {
			btn = tb.InlineButton{
				Unique: "view_balance",
				Text:   "💳 Перевірити баланс",
				Data:   strconv.FormatInt(member.ID, 10),
			}
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s \n- Роль: %s\n- ID: %d",
			member.Firstname,
			member.Username,
			status,
			role,
			member.ID,
		)

		markup := &tb.ReplyMarkup{}
		markup.InlineKeyboard = [][]tb.InlineButton{
			{btn},
		}
		c.Send(text, markup)
	}
	return nil
}

func (h *Handler) ProcessViewBalance(c tb.Context) error {
	data := c.Callback().Data

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	btnBlack := tb.InlineButton{
		Unique: "choose_card",
		Text:   "Чорна",
		Data:   fmt.Sprintf("%s|black", data),
	}
	btnWhite := tb.InlineButton{
		Unique: "choose_card",
		Text:   "Біла",
		Data:   fmt.Sprintf("%s|white", data),
	}

	markup := &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{
		{btnBlack}, {btnWhite},
	}}

	return c.Edit("🔘 Обери тип картки:", markup)
}

func (h *Handler) ProcessChooseCard(c tb.Context) error {
	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 2 {
		return c.Send("Некоректні дані.")
	}
	memberID, cardType := parts[0], parts[1]

	currencies := []struct {
		Code string
		Name string
	}{
		{"UAH", "Гривні"},
		{"PLN", "Злоті"},
		{"USD", "Долари"},
	}

	buttons := [][]tb.InlineButton{}
	for _, cur := range currencies {
		btn := tb.InlineButton{
			Unique: "final_balance",
			Text:   cur.Name,
			Data:   fmt.Sprintf("%s|%s|%s", memberID, cardType, cur.Code),
		}
		buttons = append(buttons, []tb.InlineButton{btn})
	}

	markup := &tb.ReplyMarkup{InlineKeyboard: buttons}
	return c.Edit("💱 Обери валюту:", markup)
}

func (h *Handler) ProcessFinalBalance(c tb.Context) error {
	ctx := context.Background()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 3 {
		return c.Send("Некоректні дані.")
	}
	memberID, cardType, currency := parts[0], parts[1], parts[2]

	memberIDInt, err := strconv.ParseInt(memberID, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	balance, err := h.usecase.GetBalance(ctx, us.Family.ID, memberIDInt, cardType, currency)
	if err != nil {
		h.sl.Error("failed to get balance", slog.String("err", err.Error()))
		switch e := err.(type) {
		case *errorsx.CustomError[struct{}]:
			if e.Code == errorsx.ErrCodeTokenNotFound {
				return c.Send("В данного користувача не доданий токен для перевірки балансу.")
			}
		case *errorsx.CustomError[float64]:
			if e.Code == errorsx.ErrRequestCooldown {
				return c.Send(fmt.Sprintf("Зачекай %.0f секунд перед використанням цієї функції.", e.Data))
			}
		}
		return c.Send("Не вдалося отримати баланс.")
	}

	text := fmt.Sprintf(
		"💳 Баланс (ID: %s)\nКартка: %s\nВалюта: %s\nСума: %.2f",
		memberID, cardType, currency, balance,
	)
	return c.Edit(text)
}
