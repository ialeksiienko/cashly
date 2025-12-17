package handlers

import (
	"cashly/internal/pkg/errorsx"
	"errors"
	"fmt"
	"time"

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

func mapErrorToMessage(err error) string {
	var e errorsx.ErrorWithCode
	if !errors.As(err, &e) {
		return "Сталася невідома помилка. Спробуй пізніше."
	}

	switch e.GetCode() {
	case errorsx.ErrCodeNoPermission:
		return "У тебе немає прав на цю дію."
	case errorsx.ErrCodeFailedToGenerateInviteCode:
		return "Не вдалося створити код. Спробуй пізніше."
	case errorsx.ErrCodeCannotRemoveSelf:
		return "Ти не можеш виконати цю дію з собою."
	case errorsx.ErrCodeUserHasNoFamily:
		return "У тебе поки немає жодної сім'ї. Створи або приєднайся."
	case errorsx.ErrCodeFamilyHasNoMembers:
		return "У твоїй сім'ї поки немає учасників."
	case errorsx.ErrCodeFamilyNotFound:
		return "Сім'ю з цим кодом запрошення не знайдено."
	case errorsx.ErrCodeTokenNotFound:
		return "В данного користувача не доданий токен для перевірки балансу."
	case errorsx.ErrRequestCooldown:
		return fmt.Sprintf("Зачекай %.0f секунд перед використанням цієї функції.", e.GetData())
	case errorsx.ErrCodeFamilyCodeExpired:
		expiresAt, ok := e.GetData().(time.Time)
		if !ok {
			return "Код запрошення не дійсний."
		}

		return fmt.Sprintf("Код запрошення не дійсний, закінчився - %s (час за Гринвічем, GMT)", expiresAt.Format("02.01.2006 о 15:04"))
	default:
		return "Сталася невідома помилка. Спробуй пізніше."
	}
}
