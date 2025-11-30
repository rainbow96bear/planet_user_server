package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/models"
)

// ---------------------- Todo DTO 구조 ----------------------

// TodoItem: CalendarInfo 내부에 포함되거나, Todo 개별 조회 시 사용되는 응답 구조
type TodoItem struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	EventID   uuid.UUID  `json:"eventId"`
	Content   string     `json:"content"`
	IsDone    bool       `json:"isDone"`
	DueTime   *time.Time `json:"dueTime,omitempty"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
}

// ---------------------- Todo 요청 DTO ----------------------

// TodoUpdateStatusRequest: PATCH /me/todos/:todoId 요청 시 사용
type TodoUpdateStatusRequest struct {
	IsDone bool `json:"is_done"`
}

// ------------------------------------------------------
// Todo 변환 함수
// ------------------------------------------------------

// ToTodoDTO: models.Todos를 TodoItem DTO로 변환
func ToTodoDTO(m *models.Todos) TodoItem {
	return TodoItem{
		ID:        m.ID,
		UserID:    m.UserID,
		EventID:   m.EventID,
		Content:   m.Content,
		IsDone:    m.IsDone,
		DueTime:   m.DueTime,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

// ToTodoModel: TodoItem DTO를 models.Todos로 변환 (필요 시 확장)
func ToTodoModel(d TodoItem, userID uuid.UUID) *models.Todos {
	// 주의: 이 함수는 CalendarCreateRequest에서 사용될 수 있으나,
	// 현재 ToCalendarModelFromCreate 내에서 간소화되어 사용되고 있으므로,
	// 이 파일에 TodoItem을 models.Todos로 변환하는 함수를 추가할 수 있습니다.
	return &models.Todos{
		UserID:  userID,
		Content: d.Content,
		IsDone:  d.IsDone,
		DueTime: d.DueTime,
	}
}
