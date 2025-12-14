package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/models"
)

func ToCalendarModel(
	input model.CreateCalendarInput,
	userID uuid.UUID,
) *models.CalendarEvent {

	return &models.CalendarEvent{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       input.Title,
		Emoji:       defaultEmoji(input.Emoji),
		Description: *input.Description,

		StartAt: input.StartAt,
		EndAt:   input.EndAt,

		Visibility: "private", // 기본값 (서비스에서 덮어써도 됨)

		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func ToCalendarGraphQL(
	event *models.CalendarEvent,
) *model.Calendar {

	return &model.Calendar{
		ID:          event.ID.String(),
		Title:       event.Title,
		Emoji:       &event.Emoji,
		Description: &event.Description,

		StartAt:    event.StartAt,
		EndAt:      event.EndAt,
		Visibility: model.CalendarVisibility(event.Visibility),

		Todos:     []*model.Todo{}, // 나중에 preload or resolver에서 채우기
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

func defaultEmoji(emoji *string) string {
	if emoji == nil || *emoji == "" {
		return "📝"
	}
	return *emoji
}
