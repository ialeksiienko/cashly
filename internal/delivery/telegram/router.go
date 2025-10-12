package telegram

import (
	"encoding/json"
	"monofamily/internal/delivery/telegram/handler"
	"monofamily/internal/middleware"
	"os"

	"golang.org/x/exp/slog"
	tb "gopkg.in/telebot.v3"
)

type family struct {
	Firstname string `json:"firstname"`
	ID        int64  `json:"id"`
}

var allowedUsers map[int64]struct{}

func init() {
	content, err := os.ReadFile("family.json")
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fam := []family{}

	unmarshalErr := json.Unmarshal(content, &fam)
	if unmarshalErr != nil {
		slog.Error(unmarshalErr.Error())
		os.Exit(1)
	}

	allowedUsers = make(map[int64]struct{})
	for _, f := range fam {
		allowedUsers[f.ID] = struct{}{}
	}
}

func SetupRoutes(bot *tb.Bot, h *handler.Handler) {

	bot.Use(func(hf tb.HandlerFunc) tb.HandlerFunc {
		return func(c tb.Context) error {
			userID := c.Sender().ID

			if _, exists := allowedUsers[userID]; !exists {
				return c.Send("У Вас немає прав для користування ботом, зв'яжіться з адміністратором.")
			}

			return hf(c)
		}
	})

	bot.Handle("/start", h.Start)

	bot.Handle(tb.OnText, h.HandleText)

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

		bot.Handle(&tb.InlineButton{Unique: "go_home"}, h.GoHome)
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

			familyMenu.Handle(&handler.BtnMemberDeleteNo, h.CancelMemberDeletion)
			familyMenu.Handle(&tb.InlineButton{Unique: "delete_member_yes"}, h.ProcessMemberDeletion)
		}
	}
}
