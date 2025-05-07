package dto

type UserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfileDTO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarSrc string `json:"avatar_src"`
	HideEmail bool   `json:"hide_email"`
	IsAdmin   bool   `json:"is_admin"`
}
