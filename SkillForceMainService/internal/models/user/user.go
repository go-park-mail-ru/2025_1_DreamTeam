package usermodels

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	Salt      []byte
	HideEmail bool
}

type UserProfile struct {
	Id        int
	Name      string
	Email     string
	Bio       string
	AvatarSrc string
	HideEmail bool
	IsAdmin   bool
}
