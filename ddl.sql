CREATE TABLE usertable (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT UNIQUE CHECK (LENGTH(email) BETWEEN 6 AND 320) NOT NULL,
    password TEXT NOT NULL,
    salt TEXT CHECK (LENGTH(salt) BETWEEN 8 AND 64) NOT NULL,
    name TEXT CHECK (LENGTH(name) BETWEEN 2 AND 32),
    bio TEXT CHECK (LENGTH(bio) <= 200),
    avatar_src TEXT DEFAULT '*path_to_default*',
    hide_email BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES "usertable"(id) ON DELETE CASCADE NOT NULL,
    token TEXT NOT NULL,
    expires TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE course (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    creator_user_id BIGINT REFERENCES "usertable"(id) ON DELETE CASCADE NOT NULL,
    title TEXT CHECK (LENGTH(title) BETWEEN 8 AND 32) NOT NULL,
    description TEXT CHECK (LENGTH(description) <= 200),
    avatar_src TEXT DEFAULT '*path_to_default*',
    price INT NOT NULL,
    time_to_pass INT CHECK (time_to_pass > 0 AND time_to_pass < 10000) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- CREATE TABLE course (
--     id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
--     creator_user_id BIGINT REFERENCES "usertable"(id) ON DELETE CASCADE NOT NULL,
--     title TEXT CHECK (LENGTH(title) BETWEEN 8 AND 32) NOT NULL,
--     description TEXT CHECK (LENGTH(description) <= 200),
--     avatar_src TEXT DEFAULT '*path_to_default*',
--     price INT NOT NULL,
--     tags TEXT NOT NULL,
--     time_to_pass INT CHECK (time_to_pass > 0 AND time_to_pass < 10000) NOT NULL,
--     purchase_amount INT NOT NULL,
--     metrik INT NOT NULL,
--     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
--     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
-- );
