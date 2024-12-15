
CREATE TABLE IF NOT EXISTS  users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    category TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
     content TEXT NOT NULL,
     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
     user_id INTEGER NOT NULL,
     post_id INTEGER NOT NULL,
     FOREIGN KEY (user_id) REFERENCES users (id)  on delete cascade
     FOREIGN KEY (post_id) REFERENCES posts (id)  on delete cascade
);

CREATE TABLE IF NOT EXISTS sessions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    session TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id)  on delete cascade

);


CREATE TABLE IF NOT EXISTS LikeOrDislike (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    LikeOrDislike TEXT NOT NULL,
    user_id INTEGER NOT NULL,
    post_id INTEGER NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users (id) on delete cascade
    FOREIGN KEY (post_id) REFERENCES posts (id) on delete cascade
)