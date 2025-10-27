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

	text := strings.TrimSpace(c.Text())

	session.ClearTextState(userID)

	switch state {
	case session.StateWaitingFamilyName:
		return h.processFamilyCreation(c, text)

	case session.StateWaitingFamilyCode:
		return h.processFamilyJoin(c, strings.ToUpper(text))

	case session.StateWaitingBankToken:
		return h.processUserBankToken(c)

	case session.StateWaitingPassword:
		if c.Text() != AuthPassword {
			return c.Send("❌ Невірний пароль. Спробуй ще раз.")
		}
		return h.handlePassword(c)

	default:
		return h.handleRegularText(c)
	}
}

func (h *Handler) handleRegularText(c tb.Context) error {
	userID := c.Sender().ID

	if c.Text() == AuthPassword {
		return h.handlePassword(c)
	}

	if t, ok := LastAuthTime[userID]; !ok || time.Since(t) > AuthTimeout {
		session.SetTextState(userID, session.StateWaitingPassword)
		return c.Send("🔐 Введи пароль для доступу:")
	}

	h.sl.Warn("unexpected state in HandleText", slog.Int64("user_id", userID))

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

	c.Delete()

	authMu.Lock()
	LastAuthTime[userID] = time.Now()
	authMu.Unlock()

	if _, ok := session.GetUserState(userID); !ok {
		return h.Start(c)
	}

	return nil
}
