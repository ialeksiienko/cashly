package userservice

import (
	"context"
	"errors"
	"fmt"
	"monofamily/internal/entity"
	"monofamily/internal/errorsx"
	"net/http"
	"sync"
	"time"
)

var currencyCodes = map[string]int32{
	"PLN": 985,
	"USD": 840,
}

type Currency struct {
	CurrencyCodeA int32   `json:"currencyCodeA"`
	CurrencyCodeB int32   `json:"currencyCodeB"`
	Date          int64   `json:"date"`
	RateSell      float32 `json:"rateSell"`
	RateBuy       float32 `json:"rateBuy"`
	RateCross     float32 `json:"rateCross"`
}

type Client struct {
	ClientID string    `json:"clientId"`
	Name     string    `json:"name"`
	Accounts []Account `json:"accounts"`
}

type Account struct {
	Balance int64  `json:"balance"`
	Type    string `json:"type"`
}

type tokenProvider interface {
	Get(ctx context.Context, familyID int, userID int64) (bool, *entity.UserBankToken, error)
}

func (s *UserService) GetBalance(ctx context.Context, familyID int, userID int64, cardType string, currency string) (float64, error) {

	hasToken, ubt, err := s.tokenProvider.Get(ctx, familyID, userID)
	if err != nil {
		return 0, err
	}

	if !hasToken {
		return 0, errorsx.New("user did not add bank token to family", errorsx.ErrCodeTokenNotFound, struct{}{})
	}

	currencies := []Currency{}

	curReqErr := s.handleRequest(userID, http.MethodGet, s.monoApiUrl+"bank/currency", ubt.Token, &currencies)
	if curReqErr != nil {
		return 0, curReqErr
	}

	client := Client{}

	clReqErr := s.handleRequest(userID, http.MethodGet, s.monoApiUrl+"personal/client-info", ubt.Token, &client)
	if clReqErr != nil {
		return 0, clReqErr
	}

	var balance float64
	for _, acc := range client.Accounts {
		if acc.Type == cardType {
			bal := float64(acc.Balance) / 100.0
			balance = bal
			break
		}
	}

	balance, err = convertBalance(balance, currency, currencies)
	if err != nil {
		return 0, err
	}

	return balance, nil
}

func convertBalance(balance float64, currency string, currencies []Currency) (float64, error) {
	if currency == "UAH" {
		return balance, nil
	}
	code, ok := currencyCodes[currency]
	if !ok {
		return 0, errors.New("cannot find a currency code in the map")
	}

	var rate float64

	for _, c := range currencies {
		if c.CurrencyCodeA == code && c.CurrencyCodeB == 980 {
			if code == 985 { // PLN
				rate = float64(c.RateCross)
			} else { // USD or else
				rate = float64(c.RateBuy)
			}

			if rate == 0 {
				return 0, fmt.Errorf("invalid currency rate for %s (zero)", currency)
			}

			return balance / rate, nil
		}
	}

	return 0, fmt.Errorf("currency rate for %s not found", currency)
}

var (
	cooldown    = 60 * time.Second
	lastRequest = make(map[int64]time.Time)
	mu          sync.Mutex
)

func (s *UserService) handleRequest(userID int64, method, url string, token string, obj any) error {
	now := time.Now()

	mu.Lock()
	if last, ok := lastRequest[userID]; ok && now.Sub(last) < cooldown {
		return errorsx.New("api request cooldown", errorsx.ErrRequestCooldown, (cooldown - now.Sub(last)).Seconds())
	}
	lastRequest[userID] = now
	mu.Unlock()

	reqErr := s.ApiRequest(method, url, token, obj)
	if reqErr != nil {
		return reqErr
	}

	return nil
}
