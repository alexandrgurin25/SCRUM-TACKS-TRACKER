package models

import "time"

type Comment struct {
	ID          string    `db:"id" redis:"-"`
	AuthorID    string    `db:"author_id" redis:"author_id"`
	TaskID      string    `db:"task_id" redis:"-"`
	Title       string    `db:"title" redis:"title"`
	Description string    `db:"description" redis:"description"`
	CreatedAt   time.Time `db:"created_at" redis:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" redis:"updated_at"`
}
