package handlers

import (
	"errors"
	"fmt"
	"main-service/internal/entities"
	"main-service/internal/sessions"
	"main-service/internal/usecases"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

const (
	familiesPerPage = 5
)

func (h *Handler) SelectMyFamily(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data

	isAdmin, family, err := h.usecases.FamilyService.SelectFamily(userID, data)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeFamilyNotFound {
				return c.Send("Сім'ю не знайдено.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error())
	}

	sessions.SetUserState(userID, &sessions.UserState{
		Family: family,
	})

	rows := []tb.Row{
		menu.Row(MenuViewBalance),
		menu.Row(MenuViewMembers, MenuLeaveFamily),
	}
	if isAdmin {
		rows = append(rows,
			menu.Row(MenuCreateNewCode, MenuDeleteFamily),
		)
	}
	rows = append(rows, menu.Row(MenuGoHome))

	menu.Reply(rows...)

	c.Delete()

	return c.Send(fmt.Sprintf("Увійдено в сім’ю: *%s*\n\n📂 Меню сім'ї:", family.Name), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, menu)
}

func (h *Handler) NextPage(c tb.Context) error {
	userID := c.Sender().ID

	session, exists := sessions.GetUserPageState(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	session.Page++
	sessions.SetUserPageState(userID, session)

	return showFamilyListPage(c, session.Families, session.Page)
}

func (h *Handler) PrevPage(c tb.Context) error {
	userID := c.Sender().ID

	session, exists := sessions.GetUserPageState(userID)
	if !exists {
		return c.Send("Сесія не знайдена.")
	}

	session.Page--
	sessions.SetUserPageState(userID, session)

	return showFamilyListPage(c, session.Families, session.Page)
}

func showFamilyListPage(c tb.Context, families []entities.Family, page int) error {
	if page == 0 {
		return c.Send("Це вже перша сторінка.")
	}

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

	return c.Edit("Оберіть сім’ю:", &tb.ReplyMarkup{
		InlineKeyboard: keyboard,
	})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
