package coursemodels

type Course struct {
	Id          int
	Price       int
	CreatorId   int
	TimeToPass  int
	Title       string
	Description string
	ScrImage    string
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
	PartId  int
	Lessons []*LessonPoint
}

type LessonPoint struct {
	LessonId int
	Title    string
	Type     string
	IsDone   bool
}
