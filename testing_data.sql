INSERT INTO usertable (email, password, salt, name, bio, avatar_src, hide_email) VALUES
    ('user1@example.com', 'password123', 'somesalt1', 'John Doe', 'Software Engineer and tech enthusiast.', '*path_to_default*', FALSE),
    ('user2@example.com', 'securepass456', 'somesalt2', 'Jane Smith', 'Loves data science and AI.', '*path_to_default*', TRUE),
    ('user3@example.com', 'mypassword789', 'somesalt3', 'Alice Johnson', 'Passionate about education.', '*path_to_default*', FALSE),
    ('user4@example.com', 'testpass321', 'somesalt4', 'Bob Brown', 'Enjoys programming and gaming.', '*path_to_default*', TRUE),
    ('user5@example.com', 'randompass654', 'somesalt5', 'Charlie White', 'Aspiring entrepreneur.', '*path_to_default*', FALSE);

INSERT INTO course (creator_user_id, title, description, price, time_to_pass) VALUES
    (1, 'Intro to Programming', 'Learn the basics of programming.', 100, 30),
    (2, 'Data Science Fundamentals', 'A beginner-friendly introduction to data science.', 150, 45),
    (3, 'Web Development Bootcamp', 'Become a full-stack developer.', 200, 60),
    (4, 'Machine Learning Basics', 'Understand the fundamentals of ML.', 180, 50),
    (5, 'Cybersecurity Essentials', 'Protect yourself online.', 120, 40),
    (1, 'Cloud Computing 101', 'Introduction to cloud computing concepts.', 130, 35),
    (2, 'Artificial Intelligence Basics', 'Fundamentals of AI and its applications.', 160, 55),
    (3, 'Python for Beginners', 'Learn Python programming from scratch.', 90, 25),
    (4, 'Blockchain and Cryptocurrency', 'Understand blockchain technology.', 175, 45),
    (5, 'Mobile App Development', 'Create your own mobile apps.', 190, 50),
    (1, 'Graphic Design Fundamentals', 'Learn the basics of graphic design.', 110, 30),
    (2, 'Digital Marketing Masterclass', 'Become a pro in digital marketing.', 140, 40),
    (3, 'Game Development with Unity', 'Create amazing games with Unity.', 210, 65),
    (4, 'Networking Essentials', 'Learn the core concepts of networking.', 125, 35),
    (5, 'IT Support and Helpdesk', 'Start a career in IT support.', 100, 30),
    (1, 'Ethical Hacking', 'Learn the principles of ethical hacking.', 180, 50),
    (2, 'Big Data Analytics', 'Introduction to big data technologies.', 170, 55),
    (3, 'Frontend Development', 'Master frontend technologies.', 150, 45),
    (4, 'Backend Development', 'Build scalable backend systems.', 160, 50),
    (5, 'DevOps Fundamentals', 'Learn DevOps practices and tools.', 200, 60);
