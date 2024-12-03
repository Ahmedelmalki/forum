-- Create Users Table
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

-- Create Posts Table
CREATE TABLE posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create Comments Table
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    post_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- Create Categories Table
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

-- Create Post Categories Table
CREATE TABLE post_categories (
    post_id INTEGER NOT NULL,
    category_id INTEGER NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (category_id) REFERENCES categories (id),
    PRIMARY KEY (post_id, category_id)
);

-- Create Likes Table
CREATE TABLE likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    post_id INTEGER,
    comment_id INTEGER,
    is_like BOOLEAN NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (post_id) REFERENCES posts (id),
    FOREIGN KEY (comment_id) REFERENCES comments (id)
);

-- Insert Sample Data
-- Users
INSERT INTO users (username, email, password) VALUES 
('john_doe', 'john@example.com', 'password123'),
('jane_doe', 'jane@example.com', 'password456');

-- Categories
INSERT INTO categories (name) VALUES 
('General'),
('Tech'),
('Lifestyle');

-- Posts
INSERT INTO posts (user_id, title, content) VALUES 
(1, 'Welcome to the Forum', 'This is the first post in the forum!'),
(2, 'Tech Trends 2024', 'Discuss the latest trends in technology.');

-- Post Categories
INSERT INTO post_categories (post_id, category_id) VALUES 
(1, 1), -- General
(2, 2); -- Tech

-- Comments
INSERT INTO comments (post_id, user_id, content) VALUES 
(1, 2, 'Thank you for creating this forum!'),
(2, 1, 'Great topic! Looking forward to the discussion.');

-- Likes
INSERT INTO likes (user_id, post_id, is_like) VALUES 
(1, 1, true),  -- John likes the first post
(2, 2, true);  -- Jane likes the second post
