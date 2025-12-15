package mapper

import (
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/models"
)

func ToTodoGraphQL(todo *models.Todo) *model.Todo {
	return &model.Todo{
		ID:        todo.ID.String(),
		Content:   todo.Content,
		IsDone:    todo.IsDone,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}
}

func ToCalendarGraphQL(event *models.CalendarEvent) *model.Calendar {
	if event == nil {
		return nil
	}

	todos := make([]*model.Todo, 0, len(event.Todos))
	for i := range event.Todos {
		todos = append(todos, ToTodoGraphQL(&event.Todos[i]))
	}

	return &model.Calendar{
		ID:          event.ID.String(),
		Title:       event.Title,
		Emoji:       &event.Emoji,
		Description: &event.Description,
		StartAt:     event.StartAt,
		EndAt:       event.EndAt,
		Visibility:  model.CalendarVisibility(event.Visibility),
		Todos:       todos, // ✅ 실제 데이터 전달
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
