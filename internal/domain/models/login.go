package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

type Login struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (r *Login) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Username,
			validation.Required,
			validation.Length(3, 20),
			is.Alphanumeric,
		),
		validation.Field(&r.Password,
			validation.Required,
			validation.Length(8, 64),
			validation.Match(regexp.MustCompile(`[!@#$%^&*]`)).Error("must contain a special char"),
		),
	)
}
