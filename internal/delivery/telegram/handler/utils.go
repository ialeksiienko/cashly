package handler

import (
	tb "gopkg.in/telebot.v3"
)

var generateFamilyMenu = func(isAdmin, userTokenFound bool) []tb.Row {
	rows := []tb.Row{menu.Row(MenuViewBalance)}

	if isAdmin {
		rows = append(rows, menu.Row(MenuCreateNewCode, MenuDeleteFamily))
	} else {
		rows = append(rows, menu.Row(MenuLeaveFamily))
	}

	tokenRow := menu.Row(MenuAddBankToken)
	if userTokenFound {
		tokenRow = menu.Row(MenuRemoveBankToken)
	}
	rows = append(rows, tokenRow)

	rows = append(rows,
		menu.Row(MenuViewMembers),
		menu.Row(MenuGoHome),
	)

	return rows
}
