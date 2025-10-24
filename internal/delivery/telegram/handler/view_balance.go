package handler

import (
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	tb "gopkg.in/telebot.v3"
)

type MemberID int64

var GoBackMap = make(map[int64]MemberID)

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
				return c.Send("У твоїй сім'ї поки немає учасників.")
			}
		}
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

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

		memberID, ok := GoBackMap[userID]
		if ok {
			if memberID == MemberID(member.ID) {
				c.Edit(text, markup)
				delete(GoBackMap, userID)
				break
			}
			continue
		}

		c.Send(text, markup)
	}
	return nil
}

func (h *Handler) ProcessViewBalance(c tb.Context) error {
	data := c.Callback().Data
	userID := c.Sender().ID

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	checkedUserID, err := strconv.Atoi(data)
	if err != nil {
		h.sl.Error("failed to conv user id string to int", slog.String("err", err.Error()))
		return c.Send(ErrInternalServerForUser.Error())
	}

	buttons := [][]tb.InlineButton{}

	if checkedUserID == int(userID) {
		buttons = append(buttons, []tb.InlineButton{{
			Unique: "choose_card",
			Text:   "◼️ Чорна",
			Data:   fmt.Sprintf("%s|black", data),
		}})
	}

	buttons = append(buttons, []tb.InlineButton{{Unique: "choose_card", Text: "◽️ Біла", Data: fmt.Sprintf("%s|white", data)}}, []tb.InlineButton{{Unique: "go_back", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(checkedUserID), 10)}})

	markup := &tb.ReplyMarkup{InlineKeyboard: buttons}

	return c.Edit("🔘 Обери тип картки:", markup)
}

func (h *Handler) ProcessChooseCard(c tb.Context) error {
	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 2 {
		return c.Send("Некоректні дані.")
	}
	checkedUserID, cardType := parts[0], parts[1]

	currencies := []struct {
		Code string
		Name string
	}{
		{"UAH", "₴  (Гривні)"},
		{"PLN", "zł (Злоті)"},
		{"USD", "$  (Долари)"},
	}

	buttons := [][]tb.InlineButton{}
	for _, cur := range currencies {
		btn := tb.InlineButton{
			Unique: "final_balance",
			Text:   cur.Name,
			Data:   fmt.Sprintf("%s|%s|%s", checkedUserID, cardType, cur.Code),
		}
		buttons = append(buttons, []tb.InlineButton{btn})
	}

	checkedUserIDInt, err := strconv.Atoi(checkedUserID)
	if err != nil {
		return c.Send("Не вдалося конвертувати ID особи яку перевіряєш. Спробуй ще раз.")
	}

	buttons = append(buttons, []tb.InlineButton{{Unique: "go_back", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(checkedUserIDInt), 10)}})

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
	checkedUserID, cardType, currency := parts[0], parts[1], parts[2]

	checkedUserIDInt, err := strconv.ParseInt(checkedUserID, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID користувача.")
	}

	balance, err := h.usecase.GetBalance(ctx, us.Family.ID, checkedUserIDInt, cardType, currency)
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
		checkedUserID, cardType, currency, balance,
	)
	return c.Edit(text)
}
