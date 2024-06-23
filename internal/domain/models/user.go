package models

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	PassHash []byte `json:"password_hash"`
}
