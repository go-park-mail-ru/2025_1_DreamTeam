package models

type Course struct {
	Id              int
	Price           int
	PurchasesAmount int
	CreatorId       int
	TimeToPass      int
	Title           string
	Description     string
	ScrImage        string
}
