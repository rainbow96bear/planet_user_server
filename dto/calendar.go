package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/models"
)

type TodoUpdateRequest struct {
	ID      *uuid.UUID `json:"id,omitempty"`
	Content *string    `json:"content,omitempty"`
	IsDone  *bool      `json:"isDone,omitempty"`
}

type CalendarUpdateRequest struct {
	Title       *string             `json:"title,omitempty"`
	Emoji       *string             `json:"emoji,omitempty"`
	Description *string             `json:"description,omitempty"`
	StartAt     *time.Time          `json:"startAt,omitempty"`
	EndAt       *time.Time          `json:"endAt,omitempty"`
	Visibility  *string             `json:"visibility,omitempty"`
	Todos       []TodoUpdateRequest `json:"todos,omitempty"`
}

func ToCalendarModel(
	input model.CreateCalendarInput,
	userID uuid.UUID,
) *models.CalendarEvent {

	event := &models.CalendarEvent{
		ID:          uuid.New(),
		UserID:      userID,
		Title:       input.Title,
		Emoji:       defaultEmoji(input.Emoji),
		Description: derefString(input.Description),
		StartAt:     input.StartAt,
		EndAt:       input.EndAt,
		Visibility:  derefVisibility(input.Visibility),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// ✅ Todo 변환
	if input.Todos != nil {
		event.Todos = make([]models.Todo, 0, len(input.Todos))
		for _, t := range input.Todos {
			event.Todos = append(event.Todos, models.Todo{
				ID:              uuid.New(),
				CalendarEventID: event.ID, // ⭐ 중요
				Content:         t.Content,
				IsDone:          false,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			})
		}
	}

	return event
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

func UpdateCalendarModelFromRequest(
	event *models.CalendarEvent,
	req *CalendarUpdateRequest,
) {
	if req.Title != nil {
		event.Title = *req.Title
	}

	if req.Emoji != nil {
		event.Emoji = *req.Emoji
	}

	if req.Description != nil {
		event.Description = *req.Description
	}

	if req.StartAt != nil {
		event.StartAt = *req.StartAt
	}

	if req.EndAt != nil {
		event.EndAt = *req.EndAt
	}

	if req.Visibility != nil {
		event.Visibility = *req.Visibility
	}

	// 🔹 Todos는 "전체 교체" 전략
	if req.Todos != nil {
		event.Todos = make([]models.Todo, 0, len(req.Todos))

		now := time.Now()

		for _, t := range req.Todos {
			todo := models.Todo{
				ID:              uuid.New(),
				CalendarEventID: event.ID,
				CreatedAt:       now,
				UpdatedAt:       now,
			}

			if t.Content != nil {
				todo.Content = *t.Content
			}

			if t.IsDone != nil {
				todo.IsDone = *t.IsDone
			}

			event.Todos = append(event.Todos, todo)
		}
	}

	event.UpdatedAt = time.Now()
}

func defaultEmoji(emoji *string) string {
	if emoji == nil || *emoji == "" {
		return "📝"
	}
	return *emoji
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func derefVisibility(v *model.CalendarVisibility) string {
	if v == nil {
		return "private"
	}

	switch *v {
	case model.CalendarVisibilityPublic,
		model.CalendarVisibilityFriends,
		model.CalendarVisibilityPrivate:
		return string(*v)
	default:
		return "private"
	}
}
