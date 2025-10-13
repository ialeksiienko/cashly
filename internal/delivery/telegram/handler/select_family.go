package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"monofamily/internal/entity"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

const (
	familiesPerPage = 5
)

func (h *Handler) SelectMyFamily(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data
	ctx := context.Background()

	familyID, err := strconv.Atoi(data)
	if err != nil {
		h.sl.Error("unable to convert family id string to int", slog.String("data", data))
		return c.Send(ErrInternalServerForUser.Error())
	}

	isAdmin, hasToken, family, err := h.usecase.SelectFamily(ctx, familyID, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeFamilyNotFound {
				return c.Send("Сім'ю не знайдено.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error())
	}

	session.SetUserState(userID, &session.UserState{
		Family: family,
	})

	rows := generateFamilyMenu(isAdmin, hasToken)

	menu.Reply(rows...)

	c.Delete()

	return c.Send(fmt.Sprintf("Увійдено в сім’ю: *%s*\n\n📂 Меню сім'ї:", family.Name), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, menu)
}

func (h *Handler) NextPage(c tb.Context) error {
	userID := c.Sender().ID

	s, exists := session.GetUserPageState(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	s.Page++
	session.SetUserPageState(userID, s)

	return showFamilyListPage(c, s.Families, s.Page)
}

func (h *Handler) PrevPage(c tb.Context) error {
	userID := c.Sender().ID

	s, exists := session.GetUserPageState(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	s.Page--
	session.SetUserPageState(userID, s)

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
