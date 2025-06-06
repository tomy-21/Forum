package models

import "time"

type User struct {
	ID                int64      `json:"id"`
	Username          string     `json:"username"`
	Email             string     `json:"email"`
	PasswordHash      string     `json:"-"`
	Role              string     `json:"role"`
	IsBanned          bool       `json:"is_banned"`
	CreatedAt         time.Time  `json:"created_at"`
	ProfilePictureURL *string    `json:"profile_picture_url,omitempty"`
	Bio               *string    `json:"bio,omitempty"`
	LastLoginAt       *time.Time `json:"last_login_at,omitempty"`
	MessageCount      *int       `json:"message_count,omitempty"`
	ThreadCount       *int       `json:"thread_count,omitempty"`
}

type Thread struct {
	ID           int64     `json:"id"`
	UserID       int64     `json:"user_id"`
	Author       *User     `json:"author,omitempty"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	State        string    `json:"state"`
	Tags         []*Tag    `json:"tags,omitempty"`
	MessageCount int       `json:"message_count,omitempty"`
}

type Message struct {
	ID        int64     `json:"id"`
	ThreadID  int64     `json:"thread_id"`
	UserID    int64     `json:"user_id"`
	Author    *User     `json:"author,omitempty"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	ImageURL  *string   `json:"image_url,omitempty"`
	Likes     int       `json:"likes"`
	Dislikes  int       `json:"dislikes"`
	Score     int       `json:"score"`
	UserVote  *int      `json:"user_vote,omitempty"`
}

type Tag struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type MessageVote struct {
	ID        int64 `json:"id"`
	MessageID int64 `json:"message_id"`
	UserID    int64 `json:"user_id"`
	VoteType  int   `json:"vote_type"`
}

type Friendship struct {
	UserOneID    int64      `json:"user_one_id"`
	UserTwoID    int64      `json:"user_two_id"`
	Status       string     `json:"status"`
	ActionUserID int64      `json:"action_user_id"`
	CreatedAt    time.Time  `json:"created_at"`
	AcceptedAt   *time.Time `json:"accepted_at,omitempty"`
}

type LoginPayload struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type RegistrationPayload struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
