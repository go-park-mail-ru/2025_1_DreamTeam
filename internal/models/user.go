package models

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     []byte `json:"-"`
}
