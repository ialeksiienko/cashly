package handler

import (
	"cashly/internal/session"
	"log/slog"
	"strings"
	"sync"
	"time"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) HandleText(c tb.Context) error {
	userID := c.Sender().ID
	state := session.GetTextState(userID)

	if state == session.StateNone {
		h.sl.Warn("unexpected state in HandleText", slog.Int64("user_id", userID), slog.Int("state", int(state)))
		return h.handleRegularText(c)
	}

	text := strings.TrimSpace(c.Text())

	switch state {
	case session.StateWaitingFamilyName:
		session.ClearTextState(userID)
		return h.processFamilyCreation(c, text)

	case session.StateWaitingFamilyCode:
		session.ClearTextState(userID)
		return h.processFamilyJoin(c, strings.ToUpper(text))

	case session.StateWaitingBankToken:
		session.ClearTextState(userID)
		return h.processUserBankToken(c)

	case session.StateWaitingPassword:
		return h.handlePassword(c)

	default:
		return h.handleRegularText(c)
	}
}

func (h *Handler) handleRegularText(c tb.Context) error {
	userID := c.Sender().ID

	if t, ok := LastAuthTime[userID]; !ok || time.Since(t) > AuthTimeout {
		session.SetTextState(userID, session.StateWaitingPassword)
		return c.Send("🔐 Введи пароль для доступу:")
	}

	return c.Send("Будь ласка, скористайся кнопками для взаємодії з ботом.")
}

var (
	LastAuthTime = make(map[int64]time.Time)
	AuthPassword = ""
	AuthTimeout  = 5 * time.Minute
	authMu       sync.Mutex
)

func (h *Handler) handlePassword(c tb.Context) error {
	userID := c.Sender().ID

	if c.Text() == AuthPassword {
		c.Delete()

		authMu.Lock()
		LastAuthTime[userID] = time.Now()
		authMu.Unlock()

		session.ClearTextState(c.Sender().ID)

		if _, ok := session.GetUserState(userID); !ok {
			return h.Start(c)
		}

		return c.Send("✅ Доступ дозволено. Можеш продовжити роботу.")
	}
	return c.Send("❌ Невірний пароль. Спробуй ще раз.")
}
