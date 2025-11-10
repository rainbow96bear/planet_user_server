package dto

import (
	"time"

	"github.com/rainbow96bear/planet_utils/model"
)

// CalendarInfo는 프론트엔드에 반환할 일정 정보
type CalendarInfo struct {
	EventID     uint64     `json:"eventId"`  // calendar_events.id
	UserUUID    string     `json:"userUuid"` // calendar_events.user_id
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Emoji       string     `json:"emoji"`
	StartAt     string     `json:"startAt"` // YYYY-MM-DD
	EndAt       string     `json:"endAt"`   // YYYY-MM-DD
	Visibility  string     `json:"visibility"`
	ImageURL    string     `json:"imageUrl,omitempty"`
	Todos       []TodoItem `json:"todos"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

// CalendarCreateRequest DTO
type CalendarCreateRequest struct {
	Title       string     `json:"title" form:"title" binding:"required"`
	Description string     `json:"description" form:"description"`
	Emoji       string     `json:"emoji" form:"emoji"`
	StartAt     string     `json:"startAt" form:"startAt" binding:"required"`
	EndAt       string     `json:"endAt" form:"endAt" binding:"required"`
	Visibility  string     `json:"visibility" form:"visibility" binding:"oneof=public friends private"`
	Todos       []TodoItem `json:"todos"`
}

// CalendarUpdateRequest DTO
type CalendarUpdateRequest struct {
	Title       string     `json:"title" form:"title"`
	Description string     `json:"description" form:"description"`
	Emoji       string     `json:"emoji" form:"emoji"`
	StartAt     string     `json:"startAt" form:"startAt"`
	EndAt       string     `json:"endAt" form:"endAt"`
	Visibility  string     `json:"visibility" form:"visibility" binding:"oneof=public friends private"`
	Todos       []TodoItem `json:"todos"`
}

// TodoItem DTO
type TodoItem struct {
	Text      string `json:"text" form:"text"`
	Completed bool   `json:"completed" form:"completed"`
}

// ---------------------- 변환 함수 ----------------------

// model.Calendar → CalendarInfo DTO
func ToCalendarInfo(cal *model.Calendar) *CalendarInfo {
	todos := make([]TodoItem, len(cal.Todos))
	for i, t := range cal.Todos {
		todos[i] = TodoItem{
			Text:      t.Content,
			Completed: t.Done,
		}
	}

	return &CalendarInfo{
		EventID:     cal.ID,
		UserUUID:    cal.UserUUID,
		Title:       cal.Title,
		Description: cal.Description,
		Emoji:       cal.Emoji,
		StartAt:     cal.StartAt.Format("2006-01-02"),
		EndAt:       cal.EndAt.Format("2006-01-02"),
		Visibility:  cal.Visibility,
		ImageURL:    derefString(cal.ImageURL),
		Todos:       todos,
		CreatedAt:   cal.CreatedAt,
		UpdatedAt:   cal.UpdatedAt,
	}
}

// model.Calendar 리스트 → CalendarInfo 리스트 변환
func ToCalendarInfoList(models []*model.Calendar) []*CalendarInfo {
	result := make([]*CalendarInfo, len(models))
	for i, m := range models {
		result[i] = ToCalendarInfo(m)
	}
	return result
}

// ---------------------- 보조 함수 ----------------------
func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
