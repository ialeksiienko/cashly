package handlers

import (
	"cashly/internal/entity"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/telebot.v3"
)

type MemberID int64

var (
	GoBackMap = make(map[int64]MemberID)
	GoBackMu  sync.RWMutex
)

func (h *Handler) ViewBalance(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	mms, err := h.usecase.GetFamilyMembers(ctx, f, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	GoBackMu.RLock()
	mid, ok := GoBackMap[uid]
	GoBackMu.RUnlock()

	for _, m := range mms {
		role := "Учасник"
		if m.IsAdmin {
			role = "Адміністратор"
		}

		btn := tb.InlineButton{}
		status := ""
		if !m.HasToken {
			status = " (користувач не додав токен)"
		} else {
			btn = tb.InlineButton{
				Unique: "view_balance",
				Text:   "💳 Перевірити баланс",
				Data:   strconv.FormatInt(m.ID, 10),
			}
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s \n- Роль: %s\n- ID: %d",
			m.Firstname,
			m.Username,
			status,
			role,
			m.ID,
		)

		markup := &tb.ReplyMarkup{}
		markup.InlineKeyboard = [][]tb.InlineButton{
			{btn},
		}

		if ok {
			if mid == MemberID(m.ID) {
				GoBackMu.Lock()
				c.Edit(text, markup)
				GoBackMu.Unlock()
				delete(GoBackMap, uid)
				break
			}
			continue
		}

		c.Send(text, markup)
	}
	return nil
}

func (h *Handler) ProcessViewBalance(c tb.Context) error {
	uid := c.Sender().ID
	d := c.Callback().Data

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	checkedUID, err := strconv.Atoi(d)
	if err != nil {
		h.logger.Error("failed to conv user id string to int", slog.String("err", err.Error()))
		return c.Send(ErrInternalServerForUser.Error())
	}

	b := [][]tb.InlineButton{}

	if checkedUID == int(uid) {
		b = append(b, []tb.InlineButton{{
			Unique: "choose_card",
			Text:   "◼️ Чорна",
			Data:   fmt.Sprintf("%s|black", d),
		}})
	}

	b = append(b, []tb.InlineButton{{Unique: "choose_card", Text: "◽️ Біла", Data: fmt.Sprintf("%s|white", d)}}, []tb.InlineButton{{Unique: "go_back", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(checkedUID), 10)}})

	markup := &tb.ReplyMarkup{InlineKeyboard: b}

	return c.Edit("🔘 Обери тип картки:", markup)
}

func (h *Handler) ProcessChooseCard(c tb.Context) error {
	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 2 {
		return c.Send("Некоректні дані.")
	}

	checkedUID, cardType := parts[0], parts[1]

	currencies := []struct {
		Code string
		Name string
	}{
		{"UAH", "₴  (Гривні)"},
		{"PLN", "zł (Злоті)"},
		{"USD", "$  (Долари)"},
	}

	b := [][]tb.InlineButton{}
	for _, cur := range currencies {
		btn := tb.InlineButton{
			Unique: "final_balance",
			Text:   cur.Name,
			Data:   fmt.Sprintf("%s|%s|%s", checkedUID, cardType, cur.Code),
		}
		b = append(b, []tb.InlineButton{btn})
	}

	checkedUIDInt, err := strconv.Atoi(checkedUID)
	if err != nil {
		return c.Send("Не вдалося конвертувати ID особи яку перевіряєш. Спробуй ще раз.")
	}

	b = append(b, []tb.InlineButton{{Unique: "go_back", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(checkedUIDInt), 10)}})

	markup := &tb.ReplyMarkup{InlineKeyboard: b}
	return c.Edit("💱 Обери валюту:", markup)
}

func (h *Handler) ProcessFinalBalance(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	parts := strings.Split(c.Callback().Data, "|")
	if len(parts) != 3 {
		return c.Send("Некоректні дані.")
	}

	checkedUID, cardType, currency := parts[0], parts[1], parts[2]

	checkedUIDInt, err := strconv.ParseInt(checkedUID, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID користувача.")
	}

	balance, err := h.usecase.GetBalance(ctx, f.ID, checkedUIDInt, cardType, currency)
	if err != nil {
		h.logger.Error("failed to get balance", slog.String("err", err.Error()))
		return c.Send(mapErrorToMessage(err))
	}

	if checkedUIDInt != uid {
		h.eventCh <- entity.EventNotification{
			Type:        entity.EventBalanceChecked,
			RecipientID: checkedUIDInt,
			FamilyName:  f.Name,
			Data: map[string]any{
				"checked_by_user_id": uid,
			},
		}
	}

	text := fmt.Sprintf(
		"💳 Баланс (ID: %s)\nКартка: %s\nВалюта: %s\nСума: %.2f",
		checkedUID, cardType, currency, balance,
	)
	return c.Edit(text, &tb.ReplyMarkup{InlineKeyboard: [][]tb.InlineButton{{{Unique: "go_back", Text: "⬅️ Назад", Data: strconv.FormatInt(int64(checkedUIDInt), 10)}}}})
}
