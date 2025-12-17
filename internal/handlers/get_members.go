package handlers

import (
	"cashly/internal/entity"
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	tb "gopkg.in/telebot.v3"
)

var (
	DeleteMMap = make(map[int64]MemberID)
	DeleteMMu  sync.RWMutex
)

func (h *Handler) GetMembers(c tb.Context) error {
	uid := c.Sender().ID

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	f, ok := c.Get(UfsKey).(*entity.Family)
	if !ok || f == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	mms, err := h.usecase.GetFamilyMembers(ctx, f, uid)
	if err != nil {
		return c.Send(mapErrorToMessage(err))
	}

	DeleteMMu.RLock()
	mid, ok := DeleteMMap[uid]
	DeleteMMu.RUnlock()

	for _, m := range mms {
		role := "Учасник"
		if m.IsAdmin {
			role = "Адміністратор"
		}

		userLabel := ""
		if m.IsCurrent {
			userLabel = " (це ви)"
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s\n- Роль: %s\n- ID: %d",
			m.Firstname,
			m.Username,
			userLabel,
			role,
			m.ID,
		)

		isAdmin := uid == f.CreatedBy

		if !m.IsCurrent && isAdmin {
			btn := tb.InlineButton{
				Unique: "delete_member",
				Text:   "🗑 Видалити",
				Data:   strconv.FormatInt(m.ID, 10),
			}

			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{
				{btn},
			}

			if ok {
				DeleteMMu.Lock()
				delete(DeleteMMap, uid)
				DeleteMMu.Unlock()

				if mid == MemberID(m.ID) {
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
		return c.Send(fmt.Sprintf("Всього учасників: %d", len(mms)))
	}

	return nil
}
