CREATE USER skillforce_app_user_service WITH PASSWORD 'password';

GRANT CONNECT ON DATABASE postgres TO skillforce_app_user_service;

GRANT USAGE ON SCHEMA public TO skillforce_app_user_service;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE 
    favourite_courses,
    sended_mails,
    sessions,
    signups,
    usertable
TO skillforce_app_user_service;
