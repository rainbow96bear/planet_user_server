package mapper

import (
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/models"
)

func ToCalendarGraphQL(event *models.CalendarEvent) *model.Calendar {
	if event == nil {
		return nil
	}

	return &model.Calendar{
		ID:          event.ID.String(),
		Title:       event.Title,
		Emoji:       &event.Emoji,
		Description: &event.Description,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
		Visibility:  model.CalendarVisibility(event.Visibility),
		Todos:       []*model.Todo{}, // 필요 시 resolver에서 채움
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
	}
}

func ToCalendarGraphQLList(events []*models.CalendarEvent) []*model.Calendar {
	result := make([]*model.Calendar, 0, len(events))
	for _, e := range events {
		result = append(result, ToCalendarGraphQL(e))
	}
	return result
}
