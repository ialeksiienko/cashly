package handlers

import (
	"cashly/internal/state"
	tb "gopkg.in/telebot.v3"
	"strings"
)

func (h *Handler) HandleText(c tb.Context) error {
	uid := c.Sender().ID
	s := state.GetTextState(uid)
	text := strings.TrimSpace(c.Text())

	state.ClearTextState(uid)

	switch s {
	case state.WaitingFamilyName:
		return h.processFamilyCreation(c, text)

	case state.WaitingFamilyCode:
		return h.processFamilyJoin(c, strings.ToUpper(text))

	case state.WaitingBankToken:
		return h.processUserBankToken(c)

	default:
		return nil
	}
}
