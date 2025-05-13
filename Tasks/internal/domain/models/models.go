package models

import "time"

type Task struct {
	ID          string    `redis:"-"`
	Title       string    `redis:"title"`
	Description string    `redis:"description"`
	AuthorID    string    `redis:"author_id"`
	BoardID     string    `redis:"-"`
	CreatedAt   time.Time `redis:"created_at"`
	UpdatedAt   time.Time `redis:"updated_at"`
	Deadline    time.Time `redis:"deadline"`
}
