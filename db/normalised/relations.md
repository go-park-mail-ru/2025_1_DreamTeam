# Функциональные зависимости таблицы USER
**ID → Email, Password, Salt, Name, Bio, hide_email, Avatar_src, Created_at, Updated_at**  

# Функциональные зависимости таблицы SESSIONS (сессии пользователя)
**ID → User_ID, Token, Expire, Created_at, Updated_at** 

# Функциональные зависимости таблицы COURSE
**ID → Creator_User_ID, Title, Description, Avatar_src, Price, Time_to_pass, Created_at, Updated_at** 

# Функциональные зависимости таблицы TAGS (связующая в many-to-many между VALID_TAGS и COURSE)
**ID → Tag_ID, Course_ID, Created_at, Updated_at** 

# Функциональные зависимости таблицы VALID_TAGS (есть только ограниченное количество тегов, которые можно выбрать)
**ID → Title, Created_at, Updated_at**

# Функциональные зависимости таблицы COURSE_METRIC
**ID → Course_ID, User_ID, Rating, Created_at, Updated_at**

# Функциональные зависимости таблицы PART (курс состоит из нескольких глав)
**ID → Course_ID, Order, Created_at, Updated_at**

# Функциональные зависимости таблицы LESSON_BUCKET (глава состоиз из одного или нескольких уроков)
**ID → Part_ID, Order, Created_at, Updated_at**

# Функциональные зависимости таблицы LESSON (урок состоит из нескольких блоков)
**ID → Lesson_Bucket_ID, Order, Created_at, Updated_at**

# Функциональные зависимости таблицы LESSON_CHECKPOINT (таблица с записями, которые отображают, какеи уроки прошёл пользователь)
**ID → User_ID, Lesson_ID, Created_at, Updated_at**

# Функциональные зависимости таблицы DISCUSSION (после верного решения задачи у пользователя открывается раздел обсуждений)
**ID → Lesson_ID, Author_AD, Title, Body, Created_at, Updated_at**

# Функциональные зависимости таблицы COMMENTS (пока пользователь не решил задачу, ему доступен только раздел с комментариями)
**ID → Lesson_ID, Author_AD, Title, Body, Created_at, Updated_at**

# Функциональные зависимости таблицы TEXT_LESSON (один из возможных уроков - текстовый)
**ID → Lesson_ID, Title, Created_at, Updated_at**

# Функциональные зависимости таблицы TEXT_LESSON_BLOCK (текстовый урок состоит из нескольких или одного текстового блока; тут же вместо текста в поле Value может храниться ссылка на картинку; так сделано для удобства в автоматизации создания курсов)
**ID → Text_lesson_ID, Value, Is_Image, Order, Created_at, Updated_at**

# Функциональные зависимости таблицы VIDEO_LESSON (один из возможных уроков - урок с видео)
**ID → Lesson_ID, Title, Video_src, Created_at, Updated_at**

# Функциональные зависимости таблицы TEST_LESSON (один из возможных уроков - проверочный)
**ID → Lesson_ID, Title, Created_at, Updated_at**

# Функциональные зависимости таблицы QUIZ_TASK (один из возможных вариантов проверки знаний - тест)
**ID → Lesson_Test_ID, Question, Created_at, Updated_at**

# Функциональные зависимости таблицы ANSWER_VARIANT (варианты теста)
**ID → Quiz_Task_ID, Answer, Is_True, Created_at, Updated_at**

# Функциональные зависимости таблицы PROG_TASK (один из возможных вариантов проверки знаний - проверка работоспособности кода)
**ID → Lesson_Test_ID, Description, Input_Data, Output_Data, Created_at, Updated_at**

# Функциональные зависимости таблицы ATTEMPT (попытка прохождения проверки работоспособности кода)
**ID → Prog_Task_ID, User_ID, Src, Status, Created_at, Updated_at**

