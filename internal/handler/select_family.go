package handler

import (
	"cashly/internal/entity"
	"cashly/internal/state"
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	tb "gopkg.in/telebot.v3"
)

const (
	familiesPerPage = 5
)

func (h *Handler) SelectMyFamily(c tb.Context) error {
	uid := c.Sender().ID
	d := c.Callback().Data

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	familyID, err := strconv.Atoi(d)
	if err != nil {
		h.logger.Error("unable to convert family id string to int", slog.String("data", d))
		return c.Send(ErrInternalServerForUser.Error())
	}

	isAdmin, hasToken, f, err := h.usecase.SelectFamily(ctx, familyID, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	state.SetUserState(uid, &state.UserState{
		Family: f,
	})

	rows := generateFamilyMenu(isAdmin, hasToken)

	menu.Reply(rows...)

	c.Delete()

	return c.Send(fmt.Sprintf("Увійдено в сім’ю: *%s*\n\n📂 Меню сім'ї:", f.Name), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, menu)
}

func (h *Handler) NextPage(c tb.Context) error {
	uid := c.Sender().ID

	s, exists := state.GetUserPageState(uid)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	s.Page++
	state.SetUserPageState(uid, s)

	return showFamilyListPage(c, s.Families, s.Page)
}

func (h *Handler) PrevPage(c tb.Context) error {
	uid := c.Sender().ID

	s, ok := state.GetUserPageState(uid)
	if !ok {
		return c.Send("Сесія не знайдена.")
	}

	s.Page--
	state.SetUserPageState(uid, s)

	return showFamilyListPage(c, s.Families, s.Page)
}

func showFamilyListPage(c tb.Context, families []entity.Family, page int) error {
	start := page * familiesPerPage
	totalFamilies := len(families)

	if start >= totalFamilies {
		return c.Send("Це вже остання сторінка.")
	}

	end := min(start+familiesPerPage, totalFamilies)
	current := families[start:end]

	var keyboard [][]tb.InlineButton
	for i, fam := range current {
		famCopy := fam
		btn := tb.InlineButton{
			Unique: "select_family",
			Data:   strconv.Itoa(fam.ID),
			Text:   fmt.Sprintf("%d. %s", start+i+1, famCopy.Name),
		}

		keyboard = append(keyboard, []tb.InlineButton{btn})
	}

	var paginationRow []tb.InlineButton
	if page > 0 {
		paginationRow = append(paginationRow, BtnPrevPage)
	}
	if (page+1)*familiesPerPage < totalFamilies {
		paginationRow = append(paginationRow, BtnNextPage)
	}
	if len(paginationRow) > 0 {
		keyboard = append(keyboard, paginationRow)
	}

	keyboard = append(keyboard, []tb.InlineButton{BtnGoHome})

	return c.Edit("Обери сім’ю:", &tb.ReplyMarkup{
		InlineKeyboard: keyboard,
	})
}
