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
	return validation.ValidateStruct(r,
		validation.Field(&r.Username,
			validation.Required.Error("username is required"),
			validation.Length(3, 20).Error("must be between 3 and 20 characters"),
			is.Alphanumeric.Error("must be alphanumeric"),
		),
		validation.Field(&r.Password,
			validation.Required.Error("password is required"),
			validation.Length(8, 64).Error("must be between 8 and 64 characters"),
			validation.Match(regexp.MustCompile(`[!@#$%^&*]`)).Error("must contain a special char"),
		),
	)
}
