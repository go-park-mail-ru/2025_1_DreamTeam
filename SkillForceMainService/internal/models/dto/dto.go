//go:generate easyjson -all dto.go

package dto

//easyjson:json
type UserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//easyjson:json
type UserProfileDTO struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarSrc string `json:"avatar_src"`
	HideEmail bool   `json:"hide_email"`
	IsAdmin   bool   `json:"is_admin"`
}

//easyjson:json
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

//easyjson:json
type LessonDTO struct {
	LessonHeader LessonDtoHeader `json:"header"`
	LessonBody   LessonDtoBody   `json:"lesson_body"`
}

//easyjson:json
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

//easyjson:json
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

//easyjson:json
type CourseRoadmapDTO struct {
	Parts []*CoursePartDTO `json:"parts"`
}

//easyjson:json
type LessonPointDTO struct {
	LessonId int    `json:"lesson_id"`
	Type     string `json:"lesson_type"`
	Title    string `json:"lesson_title"`
	Value    string `json:"lesson_value"`
	IsDone   bool   `json:"is_done"`
}

//easyjson:json
type LessonBucketDTO struct {
	Id      int               `json:"bucket_id"`
	Title   string            `json:"bucket_title"`
	Lessons []*LessonPointDTO `json:"lessons"`
}

//easyjson:json
type CoursePartDTO struct {
	Id      int                `json:"part_id"`
	Title   string             `json:"part_title"`
	Buckets []*LessonBucketDTO `json:"buckets"`
}

//easyjson:json
type LessonIDRequest struct {
	Id int `json:"lesson_id"`
}

//easyjson:json
type CourseIDRequest struct {
	Id int `json:"course_id"`
}

//easyjson:json
type VideoRangeRequest struct {
	Start int64
	End   int64
}

//easyjson:json
type VideoMeta struct {
	Name string
	Size int64
}

//easyjson:json
type SurveyDTO struct {
	Questions []QuestionDTO `json:"questions"`
}

//easyjson:json
type QuestionDTO struct {
	QuestionId int    `json:"question_id"`
	Question   string `json:"question"`
	LeftLebal  string `json:"left_lebal"`
	RightLebal string `json:"right_lebal"`
	Metric     string `json:"metric"`
}

//easyjson:json
type SurveyAnswerDTO struct {
	QuestionId int `json:"question_id"`
	Answer     int `json:"answer"`
}

//easyjson:json
type SurveyMetricsDTO struct {
	Metrics []SurveyMetricDTO `json:"metrics"`
}

//easyjson:json
type SurveyMetricDTO struct {
	Type         string          `json:"type"`
	Count        int             `json:"count"`
	Avg          float64         `json:"avg"`
	Distribution []int           `json:"distribution"`
	Answers      []UserAnswerDTO `json:"answers"`
}

//easyjson:json
type UserAnswerDTO struct {
	Username string `json:"username"`
	Answer   int    `json:"answer"`
}

//easyjson:json
type QuizAnswer struct {
	AnswerID int64  `json:"answer_id"`
	Answer   string `json:"answer"`
	IsRight  bool   `json:"is_right"`
}

//easyjson:json
type UserAnswer struct {
	IsRight    bool  `json:"is_right"`
	QuestionID int64 `json:"question_id"`
	AnswerID   int64 `json:"answer_id"`
}

//easyjson:json
type Test struct {
	QuestionID int64         `json:"question_id"`
	Question   string        `json:"question"`
	Answers    []*QuizAnswer `json:"answers"`
	UserAnswer UserAnswer    `json:"user_answer"`
}

//easyjson:json
type Answer struct {
	QuestionID int `json:"question_id"`
	Answer_ID  int `json:"answer_id"`
	Course_ID  int `json:"course_id"`
}

//easyjson:json
type UserQuestionAnswer struct {
	Status string `json:"status"`
	Answer string `json:"answer"`
}

//easyjson:json
type QuestionTest struct {
	QuestionID int64              `json:"question_id"`
	Question   string             `json:"question"`
	UserAnswer UserQuestionAnswer `json:"user_answer"`
}

//easyjson:json
type AnswerQuestion struct {
	QuestionID int    `json:"question_id"`
	Answer     string `json:"answer"`
}

//easyjson:json
type RaitingItem struct {
	User   UserProfileDTO `json:"user"`
	Rating int            `json:"rating"`
}

//easyjson:json
type Raiting struct {
	Rating []RaitingItem `json:"rating"`
}

//easyjson:json
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

//easyjson:json
type CreatePaymentRequest struct {
	ReturnURL string `json:"return_url"`
	User_ID   int32
	CourseID  int32 `json:"course_id"`
}

//easyjson:json
type WebhookHandlerData struct {
	Event  string `json:"event"`
	Object struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	} `json:"object"`
}
