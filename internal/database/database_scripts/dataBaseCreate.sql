CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL
);

CREATE TABLE segments (
                        id SERIAL PRIMARY KEY,
                        slug VARCHAR(255) NOT NULL
);

CREATE TABLE user_segments (
                        user_id INT REFERENCES users(id) ON DELETE CASCADE,
                        segment_id INT REFERENCES segments(id) ON DELETE CASCADE,
                        expiration_time TIMESTAMPTZ DEFAULT NOW(),
                        PRIMARY KEY (user_id, segment_id)
);

CREATE TABLE history (
                         id SERIAL PRIMARY KEY,
                         user_id INT REFERENCES users(id) ON DELETE CASCADE,
                         segment_id INT REFERENCES segments(id),
                         operation VARCHAR(255) NOT NULL,
                         created_time TIMESTAMPTZ DEFAULT NOW()
);