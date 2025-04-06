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
	Blocks       []struct {
		Body string `json:"body"`
	} `json:"blocks"`
	Footer struct {
		NextLessonId     int `json:"next_lesson_id"`
		PreviousLessonId int `json:"previous_lesson_id"`
	} `json:"footer"`
}

type LessonDtoHeader struct {
	CourseTitle string `json:"course_title"`
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

/*
{
  "Header": {
    "CourseTitle": "Информационная безопасность",
      "Part": {
        "Order": 1,
        "Title": "Введение"
      },
      "Bucket": {
        "Order": 1,
        "Title": "Первый урок"
      },
      "Points": [
          {
            "LessonId": 0,
      "Type": <"text", "video", "test">,
      "IsDone": true
      },
     ],
  },
  "Blocks": [
      {"Body": "*html*"},
      {"Body": "*html*"},
    {"Body": "*html*"},
  ],
  "Footer": {
    "NextLessonId": 1,
    "PreviousLessonId": 0,
  },
}

*/
