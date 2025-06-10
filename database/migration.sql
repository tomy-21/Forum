CREATE DATABASE IF NOT EXISTS forum_db;
USE forum_db;

-- Roles
CREATE TABLE roles (
    role_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(16) NOT NULL
);

-- Utilisateurs
CREATE TABLE Utilisateurs (
    user_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    role_id INT NOT NULL,
    name VARCHAR(64) NOT NULL,
    email VARCHAR(64) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(role_id) ON DELETE CASCADE
);

-- CatÃ©gories
CREATE TABLE categories (
    category_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Forums
CREATE TABLE forums (
    forum_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    category_id INT NOT NULL,
    name VARCHAR(64) NOT NULL,
    description TEXT,
    FOREIGN KEY (category_id) REFERENCES categories(category_id) ON DELETE CASCADE
);

-- Topics
CREATE TABLE sujet (
    topic_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    forum_id INT NOT NULL,
    user_id INT NOT NULL,
    title VARCHAR(64) NOT NULL,
    status BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    FOREIGN KEY (forum_id) REFERENCES forums(forum_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES Utilisateurs(user_id) ON DELETE CASCADE
);

-- Messages
CREATE TABLE messages (
    message_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    topic_id INT NOT NULL,
    user_id INT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY (topic_id) REFERENCES topics(topic_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Feedbacks
CREATE TABLE reaction (
    feedback_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    message_id INT NOT NULL,
    type ENUM('like', 'dislike') NOT NULL,
    UNIQUE(user_id, message_id),
    FOREIGN KEY(user_id) REFERENCES users(user_id) ON DELETE CASCADE,
    FOREIGN KEY(message_id) REFERENCES messages(message_id) ON DELETE CASCADE
);

-- Reponses
CREATE TABLE Reponses (
    reply_id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    reply_to_id INT NOT NULL,
    content TEXT NOT NULL,
    FOREIGN KEY(reply_to_id) REFERENCES messages(message_id)
);