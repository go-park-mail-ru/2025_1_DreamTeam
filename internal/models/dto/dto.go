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
