package validate

import (
	"github.com/go-playground/validator/v10"
)

type bankToken struct {
	token string `validate:"required,ascii,len=46"`
}

func BankToken(token string) bool {
	v := validator.New()
	t := bankToken{token}
	err := v.Struct(t)
	return err == nil
}
