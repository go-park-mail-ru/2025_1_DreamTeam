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

type CourseDTO struct {
	Id              int              `json:"id"`
	Price           int              `json:"price"`
	PurchasesAmount int              `json:"purchases_amount"`
	CreatorId       int              `json:"creator_id"`
	TimeToPass      int              `json:"time_to_pass"`
	Rating          float32          `json:"rating"`
	Tags            []string         `json:"tags"`
	Title           string           `json:"title"`
	Description     string           `json:"description"`
	ScrImage        string           `json:"src_image"`
	IsPurchased     bool             `json:"is_purchased"`
	IsCompleted     bool             `json:"is_completed"`
	Parts           []*CoursePartDTO `json:"parts"`
	IsFavorite      bool             `json:"is_favorite"`
}

type LessonDTO struct {
	LessonHeader LessonDtoHeader `json:"header"`
	LessonBody   LessonDtoBody   `json:"lesson_body"`
}

type LessonDtoBody struct {
	Blocks []struct {
		Body string `json:"body"`
	} `json:"blocks"`
	Footer struct {
		NextLessonId     int `json:"next_lesson_id"`
		CurrentLessonId  int `json:"current_lesson_id"`
		PreviousLessonId int `json:"previous_lesson_id"`
	} `json:"footer"`
}

type LessonDtoHeader struct {
	CourseTitle string `json:"course_title"`
	CourseId    int    `json:"course_id"`
	Part        struct {
		Order int    `json:"order"`
		Title string `json:"title"`
	} `json:"part"`
	Bucket struct {
		Order int    `json:"order"`
		Title string `json:"title"`
	} `json:"bucket"`
	Points []struct {
		LessonId int    `json:"lesson_id"`
		Type     string `json:"type"`
		IsDone   bool   `json:"is_done"`
	}
}

type CourseRoadmapDTO struct {
	Parts []*CoursePartDTO `json:"parts"`
}

type LessonPointDTO struct {
	LessonId int    `json:"lesson_id"`
	Type     string `json:"lesson_type"`
	Title    string `json:"lesson_title"`
	Value    string `json:"lesson_value"`
	IsDone   bool   `json:"is_done"`
}

type LessonBucketDTO struct {
	Id      int               `json:"bucket_id"`
	Title   string            `json:"bucket_title"`
	Lessons []*LessonPointDTO `json:"lessons"`
}

type CoursePartDTO struct {
	Id      int                `json:"part_id"`
	Title   string             `json:"part_title"`
	Buckets []*LessonBucketDTO `json:"buckets"`
}

type LessonIDRequest struct {
	Id int `json:"lesson_id"`
}

type VideoRangeRequest struct {
	Start int64
	End   int64
}

type VideoMeta struct {
	Name string
	Size int64
}

type SurveyDTO struct {
	Questions []QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
	QuestionId int    `json:"question_id"`
	Question   string `json:"question"`
	LeftLebal  string `json:"left_lebal"`
	RightLebal string `json:"right_lebal"`
	Metric     string `json:"metric"`
}

type SurveyAnswerDTO struct {
	QuestionId int `json:"question_id"`
	Answer     int `json:"answer"`
}

type SurveyMetricsDTO struct {
	Metrics []SurveyMetricDTO `json:"metrics"`
}

type SurveyMetricDTO struct {
	Type         string          `json:"type"`
	Count        int             `json:"count"`
	Avg          float64         `json:"avg"`
	Distribution []int           `json:"distribution"`
	Answers      []UserAnswerDTO `json:"answers"`
}

type UserAnswerDTO struct {
	Username string `json:"username"`
	Answer   int    `json:"answer"`
}

type QuizAnswer struct {
	AnswerID int64  `json:"answer_id"`
	Answer   string `json:"answer"`
	IsRight  bool   `json:"is_right"`
}

type UserAnswer struct {
	IsRight    bool  `json:"is_right"`
	QuestionID int64 `json:"question_id"`
	AnswerID   int64 `json:"answer_id"`
}

type Test struct {
	QuestionID int64         `json:"question_id"`
	Question   string        `json:"question"`
	Answers    []*QuizAnswer `json:"answers"`
	UserAnswer UserAnswer    `json:"user_answer"`
}

type QuizResult struct {
	Result bool `json:"result"`
}

type UserQuestionAnswer struct {
	Status string `json:"status"`
	Answer string `json:"answer"`
}

type QuestionTest struct {
	QuestionID int64              `json:"question_id"`
	Question   string             `json:"question"`
	UserAnswer UserQuestionAnswer `json:"user_answer"`
}

type RaitingItem struct {
	User   UserProfileDTO `json:"user"`
	Rating int            `json:"rating"`
}

type Raiting struct {
	Rating []RaitingItem `json:"rating"`
}

type UserStats struct {
	Percentage            int `json:"percentage"`
	CompletedTextLessons  int `json:"completed_lessons"`
	AmountTextLessons     int `json:"amount_lessons"`
	CompletedVideoLessons int `json:"completed_videos"`
	AmountVideoLessons    int `json:"amount_videos"`
	RecievedPoints        int `json:"received_points"`
	AmountPoints          int `json:"amount_points"`
	CompletedTests        int `json:"completed_tests"`
	AmountTests           int `json:"amount_tests"`
	CompletedQuestions    int `json:"completed_questions"`
	AmountQuestions       int `json:"amount_questions"`
}
