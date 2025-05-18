CREATE USER skillforce_app_course_service WITH PASSWORD 'password';

GRANT CONNECT ON DATABASE postgres TO skillforce_app_course_service;

GRANT USAGE ON SCHEMA public TO skillforce_app_course_service;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE 
    Attempt,
    Comments,
    Course_metrik,
    Discussion,
    lesson_checkpoint,
    question_task_answers,
    user_answers
TO skillforce_app_course_service;

GRANT SELECT ON TABLE 
    answer_variant,
    course,
    lesson,
    lesson_bucket,
    part,
    prog_task,
    question_task,
    quiz_task,
    tags,
    test_lesson,
    text_lesson,
    text_lesson_block,
    valid_tags,
    video_lesson
TO skillforce_app_course_service;
