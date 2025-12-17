package router

import (
	"cashly/internal/handler"
	"cashly/internal/middleware"
	"cashly/internal/state"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func SetupRoutes(bot *tb.Bot, h *handler.Handler) {

	bot.Use(
		middleware.CheckAllowedUsers,
		middleware.Auth(h.AuthPassword),
	)

	bot.Handle(tb.OnText, h.HandleText)

	bot.Handle("/start", h.Start)

	// first buttons
	{
		bot.Handle(&handler.BtnCreateFamily, h.CreateFamily)

		bot.Handle(&handler.BtnJoinFamily, h.JoinFamily)

		bot.Handle(&handler.BtnEnterMyFamily, h.EnterMyFamily)
	}

	// enter my family
	{
		bot.Handle(&tb.InlineButton{Unique: "select_family"}, h.SelectMyFamily)

		bot.Handle(&handler.BtnNextPage, h.NextPage)

		bot.Handle(&handler.BtnPrevPage, h.PrevPage)

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
				{handler.BtnCreateFamily}, {handler.BtnJoinFamily}, {handler.BtnEnterMyFamily},
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
			familyMenu.Handle(&handler.MenuViewBalance, h.ViewBalance)

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

				handler.GoBackMu.Lock()
				handler.GoBackMap[userID] = handler.MemberID(memberID)
				handler.GoBackMu.Unlock()

				return h.ViewBalance(c)
			})
		}

		familyMenu.Handle(&handler.MenuViewMembers, h.GetMembers)

		{
			familyMenu.Handle(&handler.MenuLeaveFamily, h.LeaveFamily)

			familyMenu.Handle(&handler.BtnLeaveFamilyNo, h.CancelLeaveFamily)
			familyMenu.Handle(&handler.BtnLeaveFamilyYes, h.ProcessLeaveFamily)
		}

		familyMenu.Handle(&handler.MenuAddBankToken, h.SaveUserBankToken)

		{
			familyMenu.Handle(&handler.MenuRemoveBankToken, h.RemoveBankToken)

			familyMenu.Handle(&handler.BtnRemoveBankTokenNo, h.CancelRemoveBankToken)
			familyMenu.Handle(&handler.BtnRemoveBankTokenYes, h.ProcessRemoveBankToken)
		}

		{
			familyMenu.Handle(&handler.MenuDeleteFamily, h.DeleteFamily)

			familyMenu.Handle(&handler.BtnFamilyDeleteNo, h.CancelFamilyDeletion)
			familyMenu.Handle(&handler.BtnFamilyDeleteYes, h.ProcessFamilyDeletion)
		}

		familyMenu.Handle(&handler.MenuCreateNewCode, h.CreateNewInviteCode)

		familyMenu.Handle(&handler.MenuGoHome, h.GoHome)

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

				handler.DeleteMMu.Lock()
				handler.DeleteMMap[uid] = handler.MemberID(mid)
				handler.DeleteMMu.Unlock()

				return h.GetMembers(c)
			})
		}
	}
}
