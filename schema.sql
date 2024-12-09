-- Create Users Table
CREATE TABLE IF NOT EXISTS  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    -- note 
    password TEXT NOT NULL
);

-- Create Posts Table
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
     content TEXT NOT NULL,
     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
     user_id INTEGER NOT NULL,
     post_id INTEGER NOT NULL,
     FOREIGN KEY (user_id) REFERENCES users (id)
     FOREIGN KEY (post_id) REFERENCES posts (id)
)
