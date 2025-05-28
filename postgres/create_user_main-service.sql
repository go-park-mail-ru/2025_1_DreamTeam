CREATE USER skillforce_app_main_service WITH PASSWORD 'password';

GRANT CONNECT ON DATABASE postgres TO skillforce_app_main_service;

GRANT USAGE ON SCHEMA public TO skillforce_app_main_service;

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE 
    video_lesson
TO skillforce_app_main_service;
