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

type CourseRating struct {
	CourseId int
	Rating   float32
}

type CoursePart struct {
	Id      int
	Order   int
	Title   string
	Buckets []*LessonBucket
}

type LessonBucket struct {
	Id      int
	Order   int
	Title   string
	Lessons []*LessonPoint
}

type LessonPoint struct {
	LessonId int
	Title    string
	Type     string
	IsDone   bool
}
