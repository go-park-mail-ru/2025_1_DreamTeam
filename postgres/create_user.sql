CREATE USER skillforce_app WITH PASSWORD 'password';

GRANT CONNECT ON DATABASE postgres TO skillforce_app;

GRANT USAGE ON SCHEMA public TO skillforce_app;

GRANT SELECT ON ALL TABLES IN SCHEMA public TO skillforce_app;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE 
    Attempt,
    Comments,
    Course_metrik,
    Discussion,
    favourite_courses,
    lesson_checkpoint,
    question_task_answers,
    sended_mails,
    sessions,
    signups,
    user_answers
TO skillforce_app;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE usertable TO skillforce_app;

ALTER DEFAULT PRIVILEGES IN SCHEMA public
GRANT SELECT ON TABLES TO skillforce_app;
