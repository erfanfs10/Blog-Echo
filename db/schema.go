package db

var schema = `

CREATE TABLE IF NOT EXISTS users(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT TRUE,
    is_admin BOOLEAN DEFAULT FALSE,
    last_login DATETIME,
    verification_code VARCHAR(10),
    avatar VARCHAR(100),
    password VARCHAR(255) NOT NULL,
    created DATETIME DEFAULT NOW(),
	updated DATETIME DEFAULT NOW() ON UPDATE NOW()
);

CREATE TABLE IF NOT EXISTS posts(
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    title TEXT,
    body TEXT,
    user_id INT NOT NULL,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
    created DATETIME DEFAULT NOW(),
	updated DATETIME DEFAULT NOW() ON UPDATE NOW()
);

`
