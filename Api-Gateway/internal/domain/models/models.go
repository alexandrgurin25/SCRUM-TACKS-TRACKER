package models

import "time"

type User struct {
	ID       string
	Email    string
	Username string
}

type Task struct {
	ID          string
	Title       string
	Description string
	AuthorID    string
	BoardID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Deadline    time.Time
}

type Board struct {
	ID    string
	Title string
	Tasks []*Task
}

type Project struct {
	ID             string
	AuthorID       string
	ParticipantsID []string
	Title          string
	Description    string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Comment struct {
	ID          string
	AuthorID    string
	TaskID      string
	Title       string
	Description string
}
