package models

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey"`
	CalendarEventID uuid.UUID `gorm:"type:uuid;not null"`
	Content         string
	IsDone          bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (Todo) TableName() string {
	return "todos"
}
