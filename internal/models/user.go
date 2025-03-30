package models

type User struct {
	Id       int    `json:"-"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Salt     []byte `json:"-"`
}

type UserProfile struct {
	Id        int    `json:"-"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarSrc string `json:"avatar_src"`
	HideEmail bool   `json:"hide_email"`
}
