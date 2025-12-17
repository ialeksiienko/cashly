package router

import (
	"cashly/internal/handlers"
	"cashly/internal/middleware"
	"cashly/internal/state"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func SetupRoutes(bot *tb.Bot, h *handlers.Handler) {

	bot.Use(
		middleware.CheckAllowedUsers,
		middleware.Auth(h.AuthPassword),
	)

	bot.Handle(tb.OnText, h.HandleText)

	bot.Handle("/start", h.Start)

	// first buttons
	{
		bot.Handle(&handlers.BtnCreateFamily, h.CreateFamily)

		bot.Handle(&handlers.BtnJoinFamily, h.JoinFamily)

		bot.Handle(&handlers.BtnEnterMyFamily, h.EnterMyFamily)
	}

	// enter my family
	{
		bot.Handle(&tb.InlineButton{Unique: "select_family"}, h.SelectMyFamily)

		bot.Handle(&handlers.BtnNextPage, h.NextPage)

		bot.Handle(&handlers.BtnPrevPage, h.PrevPage)

		bot.Handle(&tb.InlineButton{Unique: "go_home"}, func(c tb.Context) error {
			userID := c.Sender().ID

			state.DeleteUserState(userID)

			{
				msg, _ := bot.Send(c.Sender(), ".", &tb.SendOptions{
					ReplyMarkup: &tb.ReplyMarkup{
						RemoveKeyboard: true,
					},
				})

				bot.Delete(msg)
			}

			inlineKeys := [][]tb.InlineButton{
				{handlers.BtnCreateFamily}, {handlers.BtnJoinFamily}, {handlers.BtnEnterMyFamily},
			}

			return c.Edit("Вибери один з варіантів на клавіатурі.", &tb.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			})
		})
	}

	familyMenu := bot.Group()
	familyMenu.Use(middleware.CheckUserState(h.GoHome))

	// family menu
	{
		{
			familyMenu.Handle(&handlers.MenuViewBalance, h.ViewBalance)

			familyMenu.Handle(&tb.InlineButton{Unique: "view_balance"}, h.ProcessViewBalance)
			familyMenu.Handle(&tb.InlineButton{Unique: "choose_card"}, h.ProcessChooseCard)
			familyMenu.Handle(&tb.InlineButton{Unique: "final_balance"}, h.ProcessFinalBalance)

			familyMenu.Handle(&tb.InlineButton{Unique: "go_back"}, func(c tb.Context) error {
				userID := c.Sender().ID
				data := c.Callback().Data

				memberID, err := strconv.Atoi(data)
				if err != nil {
					return c.Edit("Не вдалося повернутися назад.")
				}

				handlers.GoBackMu.Lock()
				handlers.GoBackMap[userID] = handlers.MemberID(memberID)
				handlers.GoBackMu.Unlock()

				return h.ViewBalance(c)
			})
		}

		familyMenu.Handle(&handlers.MenuViewMembers, h.GetMembers)

		{
			familyMenu.Handle(&handlers.MenuLeaveFamily, h.LeaveFamily)

			familyMenu.Handle(&handlers.BtnLeaveFamilyNo, h.CancelLeaveFamily)
			familyMenu.Handle(&handlers.BtnLeaveFamilyYes, h.ProcessLeaveFamily)
		}

		familyMenu.Handle(&handlers.MenuAddBankToken, h.SaveUserBankToken)

		{
			familyMenu.Handle(&handlers.MenuRemoveBankToken, h.RemoveBankToken)

			familyMenu.Handle(&handlers.BtnRemoveBankTokenNo, h.CancelRemoveBankToken)
			familyMenu.Handle(&handlers.BtnRemoveBankTokenYes, h.ProcessRemoveBankToken)
		}

		{
			familyMenu.Handle(&handlers.MenuDeleteFamily, h.DeleteFamily)

			familyMenu.Handle(&handlers.BtnFamilyDeleteNo, h.CancelFamilyDeletion)
			familyMenu.Handle(&handlers.BtnFamilyDeleteYes, h.ProcessFamilyDeletion)
		}

		familyMenu.Handle(&handlers.MenuCreateNewCode, h.CreateNewInviteCode)

		familyMenu.Handle(&handlers.MenuGoHome, h.GoHome)

		// admin menu
		{
			familyMenu.Handle(&tb.InlineButton{Unique: "delete_member"}, h.DeleteMember)

			familyMenu.Handle(&tb.InlineButton{Unique: "delete_member_no"}, h.CancelMemberDeletion)
			familyMenu.Handle(&tb.InlineButton{Unique: "delete_member_yes"}, h.ProcessMemberDeletion)

			familyMenu.Handle(&tb.InlineButton{Unique: "go_back_delete_member"}, func(c tb.Context) error {
				uid := c.Sender().ID
				d := c.Callback().Data

				mid, err := strconv.Atoi(d)
				if err != nil {
					return c.Edit("Не вдалося повернутися назад.")
				}

				handlers.DeleteMMu.Lock()
				handlers.DeleteMMap[uid] = handlers.MemberID(mid)
				handlers.DeleteMMu.Unlock()

				return h.GetMembers(c)
			})
		}
	}
}
