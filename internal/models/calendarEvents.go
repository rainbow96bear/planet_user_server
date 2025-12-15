package models

import (
	"time"

	"github.com/google/uuid"
)

// CalendarEvent represents a calendar event entity
// DB: calendar_events
type CalendarEvent struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Title       string
	Emoji       string
	Description string
	StartAt     time.Time
	EndAt       time.Time
	Visibility  string
	CreatedAt   time.Time
	UpdatedAt   time.Time

	Todos []Todo `gorm:"foreignKey:CalendarEventID;references:ID;constraint:OnDelete:CASCADE"`
}

func (CalendarEvent) TableName() string {
	return "calendar_events"
}
