package handler

import (
	"cashly/internal/errorsx"
	"cashly/internal/session"
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/telebot.v3"
)

var (
	GoBackDeleteMemberMap = make(map[int64]MemberID)
	GoBackDMMu            sync.RWMutex
)

func (h *Handler) GetMembers(c tb.Context) error {
	userID := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	members, err := h.usecase.GetFamilyMembersInfo(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeFamilyHasNoMembers {
				return c.Send("У твоїй сім'ї поки немає учасників.")
			}
		}
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

	GoBackDMMu.RLock()
	memberID, ok := GoBackDeleteMemberMap[userID]
	GoBackDMMu.RUnlock()

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

		isAdmin := userID == us.Family.CreatedBy

		if !member.IsCurrent && isAdmin {
			btn := tb.InlineButton{
				Unique: "delete_member",
				Text:   "🗑 Видалити",
				Data:   strconv.FormatInt(member.ID, 10),
			}

			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{
				{btn},
			}

			if ok {
				GoBackDMMu.Lock()
				delete(GoBackDeleteMemberMap, userID)
				GoBackDMMu.Unlock()

				if memberID == MemberID(member.ID) {
					c.Edit(text, markup)
					break
				}
				continue
			}

			c.Send(text, markup)
		} else {
			if ok {
				continue
			}
			c.Send(text)
		}
	}

	if !ok {
		return c.Send(fmt.Sprintf("Всього учасників: %d", len(members)))
	}

	return nil
}
