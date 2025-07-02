package handlers

import (
	"fmt"
	"main-service/internal/sessions"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) GetMembers(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists {
		return c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
	}

	members, err := h.usecases.UserService.GetMembersInfo(us.Family, userID)
	if err != nil {
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

	membersLen := len(members)

	if membersLen == 0 {
		return c.Send("У вашій сім'ї поки немає учасників.")
	}

	c.Send("📋 Список учасників сім'ї:\n")

	for _, member := range members {
		role := "Учасник"
		if member.IsAdmin {
			role = "Адміністратор"
		}

		userLabel := ""
		if member.IsCurrent {
			userLabel = " (це ви)"
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s\n- Роль: %s\n- ID: %d",
			member.Firstname,
			member.Username,
			userLabel,
			role,
			member.ID,
		)

		if !member.IsCurrent && member.IsAdmin {
			btn := tb.InlineButton{
				Unique: "delete_member",
				Text:   "🗑 Видалити",
				Data:   strconv.FormatInt(member.ID, 10),
			}

			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{
				{btn},
			}

			c.Send(text, markup)
		} else {
			c.Send(text)
		}
	}

	return c.Send(fmt.Sprintf("Всього учасників: %d", membersLen))
}

func (h *Handler) LeaveFamily(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists {
		return c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
	}

	if us.Family.CreatedBy == userID {
		return c.Send("Адміністратор не може вийти з сім'ї.")
	}

	err := h.usecases.UserService.LeaveFamily(us.Family.ID, userID)
	if err != nil {
		return c.Send("Не вдалося вийти з сім'ї. Спробуйте ще раз пізніше.")
	}

	sessions.DeleteUserState(userID)

	msg, _ := h.bot.Send(c.Sender(), ".", &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{
			RemoveKeyboard: true,
		},
	})

	h.bot.Delete(msg)

	inlineKeys := [][]tb.InlineButton{
		{BtnCreateFamily}, {BtnJoinFamily}, {BtnEnterMyFamily},
	}

	return c.Send(
		"Ви успішно вийшли з сім'ї.\n\nВибери один з варіантів на клавіатурі.",
		&tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		},
	)
}

// admin handlers

func (h *Handler) DeleteMember(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
	}

	if userID != us.Family.CreatedBy {
		return c.Send("У вас немає прав на видалення.")
	}

	if userID == memberID {
		return c.Send("Ви не можете видалити себе.")
	}

	removeErr := h.usecases.AdminService.RemoveMember(us.Family.ID, memberID)
	if removeErr != nil {
		return c.Send("Не вдалося видалити користувача з сім'ї. Спробуйте ще раз пізніше.")
	}

	h.bot.Edit(c.Message(), "Учасника успішно видалено. Оновлюю список...")

	return h.GetMembers(c)
}

func (h *Handler) DeleteFamily(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
	}

	if userID != us.Family.CreatedBy {
		return c.Send("У вас немає прав на видалення сім'ї.")
	}

	err := h.usecases.AdminService.DeleteFamily(us.Family.ID)
	if err != nil {
		return c.Send("Не вдалося видалити сім'ю. Спробуйте ще раз пізніше.")
	}

	sessions.DeleteUserState(userID)

	msg, _ := h.bot.Send(c.Sender(), ".", &tb.SendOptions{
		ReplyMarkup: &tb.ReplyMarkup{
			RemoveKeyboard: true,
		},
	})

	h.bot.Delete(msg)

	inlineKeys := [][]tb.InlineButton{
		{BtnCreateFamily}, {BtnJoinFamily}, {BtnEnterMyFamily},
	}

	return c.Send(
		"Сім'ю успішно видалено.\n\nВибери один з варіантів на клавіатурі.",
		&tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		},
	)
}

func (h *Handler) CreateNewInviteCode(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		return c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
	}

	if userID != us.Family.CreatedBy {
		return c.Send("У вас немає прав на створення нового коду запрошення.")
	}

	code, expiresAt, err := h.usecases.AdminService.CreateNewFamilyCode(us.Family.ID, userID)
	if err != nil {
		return c.Send("Не вдалося створити код запрошення. Спробуйте ще раз пізніше.")
	}

	return c.Send(fmt.Sprintf("Новий код запрошення: `%s`\n\nДійсний до — %s", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
}
