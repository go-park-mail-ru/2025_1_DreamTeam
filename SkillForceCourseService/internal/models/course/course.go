package coursemodels

type Course struct {
	Id          int
	Price       int
	CreatorId   int
	TimeToPass  int
	Title       string
	Description string
	ScrImage    string
	Parts       []*CoursePart
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
	Value    string
	IsDone   bool
	IsImage  bool
	BucketId int
	Order    int
}

type Survey struct {
	Id        int
	Questions []Question
}

type Question struct {
	QuestionId int
	Question   string
	LeftLebal  string
	RightLebal string
	Metric     string
}

type SurveyAnswer struct {
	Id         int
	QuestionId int
	Answer     int
}

type SurveyMetrics struct {
	Metrics []SurveyMetrics
}

type SurveyMetric struct {
	Type         string
	Count        int
	Avg          float64
	Distribution []int
	Answers      []UserAnswer
}

type UserAnswer struct {
	Username string
	Answer   int
}

type QuizAnswer struct {
	AnswerID int64
	Answer   string
}

type Test struct {
	QuestionID    int64
	Question      string
	RightAnswerID int64
	Answers       []QuizAnswer
}
