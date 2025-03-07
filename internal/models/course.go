package models

type Course struct {
	Id              int    `json:"id"`
	Price           int    `json:"price"`
	PurchasesAmount int    `json:"purchases_amount"`
	CreatorId       int    `json:"creator_id"`
	TimeToPass      int    `json:"time_to_pass"`
	Title           string `json:"title"`
	Description     string `json:"description"`
	ScrImage        string `json:"src_image"`
}
