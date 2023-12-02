CREATE TABLE files (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    content BYTEA,
    FOREIGN KEY (user_id) REFERENCES users (id)
);
