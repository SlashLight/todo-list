package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID `json:"id"`
	AuthorID    uuid.UUID `json:"author-id"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Status      string    `json:"status"`
	Deadline    time.Time `json:"deadline,omitempty"`
}
