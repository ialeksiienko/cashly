package test

import (
	userservice "cashly/internal/service/user"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConvertBalance(t *testing.T) {
	tests := []struct {
		name       string
		balance    float64
		currency   string
		currencies []userservice.Currency
		want       float64
		wantErr    bool
	}{
		{
			name:     "UAH - no conversion",
			balance:  1000.0,
			currency: "UAH",
			want:     1000.0,
		},
		{
			name:     "USD conversion",
			balance:  4150.0,
			currency: "USD",
			currencies: []userservice.Currency{
				{CurrencyCodeA: 840, CurrencyCodeB: 980, RateBuy: 41.5},
			},
			want: 100.0,
		},
		{
			name:     "PLN uses RateCross",
			balance:  985.0,
			currency: "PLN",
			currencies: []userservice.Currency{
				{CurrencyCodeA: 985, CurrencyCodeB: 980, RateCross: 9.85},
			},
			want: 100.0,
		},
		{
			name:     "unknown currency",
			currency: "XXX",
			wantErr:  true,
		},
		{
			name:     "zero rate",
			balance:  1000.0,
			currency: "USD",
			currencies: []userservice.Currency{
				{CurrencyCodeA: 840, CurrencyCodeB: 980, RateBuy: 0},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userservice.ConvertBalance(tt.balance, tt.currency, tt.currencies)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.InDelta(t, tt.want, got, 0.01)
		})
	}
}
