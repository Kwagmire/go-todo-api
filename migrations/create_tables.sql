-- Create 'users' table
CREATE TABLE IF NOT EXISTS users (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	password_hash VARCHAR(255) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create 'todos' table
CREATE TABLE IF NOT EXISTS todos (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	title VARCHAR(255) NOT NULL,
	description TEXT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);

-- Add index for user_id on todos table for faster lookups
CREATE INDEX IF NOT EXISTS idx_todos_user_id ON todos(user_id);
