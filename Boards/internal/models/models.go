package models

import "time"

type Board struct {
	ID        string    `json:"id"`
	AuthorID  string    `json:"author_id"`
	ProjectID string    `json:"project_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BoardsList struct {
	ProjectID string  `json:"project_id"`
	Boards    []Board `json:"boards"`
}

type Task struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	AuthorID    string    `json:"author_id"`
	BoardID     string    `json:"-"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Deadline    time.Time `json:"deadline"`
}

type TasksList struct {
	BoardID string `json:"board_id"`
	Tasks   []Task `json:"tasks"`
}
