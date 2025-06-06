package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB // Global DB variable can be kept for package-internal use if necessary

// Init initializes the database connection and creates tables.
// It now returns the *sql.DB instance and an error.
func Init(dataSourceName string) (*sql.DB, error) {
	var err error
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error configuring MySQL driver: %w", err)
	}

	if err = db.Ping(); err != nil {
		db.Close() // Close if ping fails
		return nil, fmt.Errorf("error connecting to MySQL database: %w", err)
	}
	log.Println("Connected to MySQL database!")

	// Set the global DB variable if you still need it for other functions within this package.
	// However, relying on the returned 'db' instance is generally cleaner.
	DB = db

	// Create tables using the established connection
	if err := CreateTables(db); err != nil {
		db.Close() // Close DB if table creation fails
		return nil, fmt.Errorf("error creating tables: %w", err)
	}

	return db, nil
}

// createTables now accepts a *sql.DB argument.
func CreateTables(db *sql.DB) error {
	userTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) UNIQUE NOT NULL,
        email VARCHAR(255) UNIQUE NOT NULL,
        password_hash VARCHAR(128) NOT NULL,
        role VARCHAR(50) NOT NULL DEFAULT 'user',
        is_banned BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        profile_picture_url TEXT,
        bio TEXT,
        last_login_at TIMESTAMP NULL DEFAULT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	threadTableSQL := `
    CREATE TABLE IF NOT EXISTS threads (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
        title VARCHAR(255) NOT NULL,
        description TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
        state VARCHAR(50) NOT NULL DEFAULT 'ouvert',
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	messageTableSQL := `
    CREATE TABLE IF NOT EXISTS messages (
        id INT AUTO_INCREMENT PRIMARY KEY,
        thread_id INT NOT NULL,
        user_id INT NOT NULL,
        content TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        image_url TEXT,
        FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	tagTableSQL := `
    CREATE TABLE IF NOT EXISTS tags (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(100) UNIQUE NOT NULL
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	threadTagsTableSQL := `
    CREATE TABLE IF NOT EXISTS thread_tags (
        thread_id INT NOT NULL,
        tag_id INT NOT NULL,
        PRIMARY KEY (thread_id, tag_id),
        FOREIGN KEY (thread_id) REFERENCES threads(id) ON DELETE CASCADE,
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	messageVotesTableSQL := `
    CREATE TABLE IF NOT EXISTS message_votes (
        id INT AUTO_INCREMENT PRIMARY KEY,
        message_id INT NOT NULL,
        user_id INT NOT NULL,
        vote_type TINYINT NOT NULL,
        UNIQUE (message_id, user_id),
        FOREIGN KEY (message_id) REFERENCES messages(id) ON DELETE CASCADE,
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	queries := []string{
		userTableSQL,
		threadTableSQL,
		messageTableSQL,
		tagTableSQL,
		threadTagsTableSQL,
		messageVotesTableSQL,
	}

	for _, query := range queries {
		_, err := db.Exec(query) // Use the passed db instance
		if err != nil {
			return fmt.Errorf("error executing query to create table: %s, error: %w", query, err)
		}
	}
	log.Println("Tables created or already exist (MySQL).")
	return nil
}
