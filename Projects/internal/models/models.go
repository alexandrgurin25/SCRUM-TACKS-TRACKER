package models

import "time"

type Project struct {
	ID          string    `json:"id"`
	AuthorID    string    `json:"author_id"`
	Description string    `json:"description"`
	Title       string    `json:"title"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ParticipantsList struct {
	ProjectID string   `json:"project_id"`
	UsersID   []string `json:"users_id"`
}
