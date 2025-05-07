package models

type User struct {
	Email    string `json:"email"`
	PassHash []byte `json:"pass_hash"`
}
