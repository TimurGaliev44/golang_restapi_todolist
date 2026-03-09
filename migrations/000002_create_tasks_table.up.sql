CREATE TABLE IF NOT EXISTS tasks(
		id SERIAL PRIMARY KEY,
        user_id INT NOT NULL REFERENCES users(id),
		title VARCHAR(100) NOT NULL,
		description VARCHAR(200) NOT NULL,
		completed BOOL NOT NULL,
		created_at TIMESTAMP NOT NULL,
		completed_at TIMESTAMP
);