package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"regexp"
)

type RegisterRequest struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type RegisterGoodResponse struct {
	Status string `json:"status"`
}

func (r *RegisterRequest) Validate() error {
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
