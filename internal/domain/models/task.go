package models

import (
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID
	AuthorID    uuid.UUID
	Title       string
	Description string
	Status      string
	Deadline    time.Time
}
