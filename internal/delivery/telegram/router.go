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
		}

		familyMenu.Handle(&handler.MenuViewMembers, h.GetMembers)

		{
			familyMenu.Handle(&handler.MenuLeaveFamily, h.LeaveFamily)

			familyMenu.Handle(&handler.BtnLeaveFamilyNo, h.CancelLeaveFamily)
			familyMenu.Handle(&handler.BtnLeaveFamilyYes, h.ProcessLeaveFamily)
		}

		familyMenu.Handle(&handler.MenuAddBankToken, h.SaveUserBankToken)

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

	//bot.Handle("/family", func(c tb.Context) error {
	//	userID := c.Sender().ID
	//
	//	families, err := repo.GetFamiliesByUserID(userID)
	//	if err != nil {
	//		sl.Error("failed to get family", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
	//		return c.Send("Не вдалося отримати дані про сім'ю.")
	//	}
	//
	//	if len(families) == 0 {
	//		return c.Send("Ти ще не приєднаний до жодної сім’ї.")
	//	}
	//
	//	msg := "Твої сім’ї:\n"
	//	for _, f := range families {
	//		users, err := repo.GetAllUsersInFamily(&f)
	//		if err != nil {
	//			sl.Error("failed to get users", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
	//			return c.Send("Не вдалося отримати всіх користувачів сім'ї. Спробуй пізніше.")
	//		}
	//
	//		usersList := make([]string, len(users))
	//		for _, u := range users {
	//			usersList = append(usersList, u.Username)
	//		}
	//
	//		msg += fmt.Sprintf("🔸 %s (ID: %d). \nУчасники: %s\n", f.Name, f.ID, strings.Join(usersList, ", "))
	//	}
	//
	//	return c.Send(msg)
	//})
	//
	//button := tb.InlineButton{
	//			Unique: "mono_link",
	//			Text:   "Силка",
	//			URL:    "https://api.monobank.ua/",
	//		}
	//
	//		inlineKeys := [][]tb.InlineButton{
	//			{button},
	//		}
	//
	//		bot.Send(c.Sender(), "Привіт, цей бот допоможе дізнатися рахунок на карті monobank.\n\nПерейди по силці внизу та відправ свій токен в цей чат.", &tb.ReplyMarkup{
	//			InlineKeyboard: inlineKeys,
	//		})
	//
	//		bot.Handle(tb.OnText, func(c tb.Context) error {
	//			return c.Send("Данные успешно сохранены в базе данных!")
	//		}, middlewares.CheckTokenValid)
	//		return nil

	//h.bot.Handle("/balance", func(c tb.Context) error {
	//	buttonBlack := tb.InlineButton{Unique: "black", Text: "Черная"}
	//	buttonWhite := tb.InlineButton{Unique: "white", Text: "Белая"}
	//
	//	inlineKeysCardType := [][]tb.InlineButton{
	//		{buttonBlack},
	//		{buttonWhite},
	//	}
	//
	//	h.bot.Send(c.Sender(), "Напиши какого типа карточки ты хотел бы узнать баланс.", &tb.ReplyMarkup{InlineKeyboard: inlineKeysCardType})
	//
	//	h.bot.Handle(&buttonBlack, func(c tb.Context) error {
	//		buttonHryvnia := tb.InlineButton{Unique: "hryvnia", Text: "Гривны"}
	//		buttonZloty := tb.InlineButton{Unique: "zloty", Text: "Злотые"}
	//		buttonDollars := tb.InlineButton{Unique: "dollars", Text: "Доллары"}
	//
	//		inlineKeysCurrency := [][]tb.InlineButton{
	//			{buttonHryvnia},
	//			{buttonZloty},
	//			{buttonDollars},
	//		}
	//
	//		h.bot.Send(c.Sender(), "Теперь нажми на кнопку в какой валюте ты хочешь узнать баланс.", &tb.ReplyMarkup{InlineKeyboard: inlineKeysCurrency})
	//
	//		h.bot.Handle(&buttonZloty, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonZloty.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		h.bot.Handle(&buttonDollars, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonDollars.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		h.bot.Handle(&buttonHryvnia, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonHryvnia.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		return nil
	//	})
	//	return nil
	//})
}
