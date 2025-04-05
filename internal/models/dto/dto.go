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
}

type CourseDTO struct {
	Id              int      `json:"id"`
	Price           int      `json:"price"`
	PurchasesAmount int      `json:"purchases_amount"`
	CreatorId       int      `json:"creator_id"`
	TimeToPass      int      `json:"time_to_pass"`
	Rating          float32  `json:"rating"`
	Tags            []string `json:"tags"`
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	ScrImage        string   `json:"src_image"`
}
