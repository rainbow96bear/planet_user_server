package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/models"
)

// ---------------------- DTO 구조 ----------------------

type CalendarInfo struct {
	EventID     uuid.UUID  `json:"eventId"`
	UserID      string     `json:"UserID"` // UUID를 문자열로 전달
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Emoji       string     `json:"emoji"`
	StartAt     string     `json:"startAt"`
	EndAt       string     `json:"endAt"`
	Visibility  string     `json:"visibility"`
	Todos       []TodoItem `json:"todos"`
	ImageURL    *string    `json:"imageUrl,omitempty"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

type CalendarCreateRequest struct {
	Title       string     `json:"title" binding:"required,max=255"`
	Description string     `json:"description" binding:"max=65535"`
	Emoji       string     `json:"emoji" binding:"omitempty,max=10"`
	StartAt     string     `json:"startAt" binding:"required,datetime=2006-01-02"`
	EndAt       string     `json:"endAt" binding:"required,datetime=2006-01-02,gtefield=StartAt"`
	Visibility  string     `json:"visibility" binding:"required,oneof=public friends private"`
	Todos       []TodoItem `json:"todos" binding:"omitempty,dive"`
	ImageURL    *string    `json:"imageUrl"`
}

type CalendarUpdateRequest struct {
	Title       *string     `json:"title,omitempty" binding:"omitempty,max=255"`
	Description *string     `json:"description,omitempty" binding:"omitempty,max=65535"`
	Emoji       *string     `json:"emoji,omitempty" binding:"omitempty,max=10"`
	StartAt     *string     `json:"startAt,omitempty" binding:"omitempty,datetime=2006-01-02"`
	EndAt       *string     `json:"endAt,omitempty" binding:"omitempty,datetime=2006-01-02"`
	Visibility  *string     `json:"visibility,omitempty" binding:"omitempty,oneof=public friends private"`
	Todos       *[]TodoItem `json:"todos,omitempty"`
	ImageURL    *string     `json:"imageUrl,omitempty"`
}

type TodoItem struct {
	EventID   *uuid.UUID `json:"eventId,omitempty"`
	Text      string     `json:"text" binding:"required,max=255"`
	Completed bool       `json:"completed"`
}

// ---------------------- 변환 함수 ----------------------

func ToCalendarInfo(cal *models.CalendarEvents) *CalendarInfo {
	if cal == nil {
		return nil
	}

	todos := make([]TodoItem, len(cal.Todos))
	for i, t := range cal.Todos {
		todos[i] = TodoItem{
			EventID:   &t.ID,
			Text:      t.Content,
			Completed: t.IsDone,
		}
	}

	return &CalendarInfo{
		EventID:     cal.ID,
		UserID:      cal.UserID.String(),
		Title:       cal.Title,
		Description: cal.Description,
		Emoji:       cal.Emoji,
		StartAt:     formatDate(cal.StartAt),
		EndAt:       formatDate(cal.EndAt),
		Visibility:  cal.Visibility,
		Todos:       todos,
		ImageURL:    &cal.ImageURL,
		CreatedAt:   cal.CreatedAt,
		UpdatedAt:   cal.UpdatedAt,
	}
}

func ToCalendarInfoList(events []*models.CalendarEvents) []*CalendarInfo {
	result := make([]*CalendarInfo, 0, len(events))
	for _, e := range events {
		if info := ToCalendarInfo(e); info != nil {
			result = append(result, info)
		}
	}
	return result
}

func ToCalendarModelFromCreate(req *CalendarCreateRequest, userID uuid.UUID) *models.CalendarEvents {
	cal := &models.CalendarEvents{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     parseDate(req.StartAt),
		EndAt:       parseDate(req.EndAt),
		Visibility:  req.Visibility,
		ImageURL:    stringOrEmpty(req.ImageURL),
	}

	// Todos가 nil이 아니면 append로 안전하게 추가
	if len(req.Todos) > 0 {
		for _, t := range req.Todos {
			cal.Todos = append(cal.Todos, &models.Todos{
				UserID:  userID,
				Content: t.Text,
				IsDone:  t.Completed,
				// EventID는 cal.ID로 나중에 GORM에서 자동 할당 가능
			})
		}
	}

	return cal
}

func UpdateCalendarModelFromRequest(cal *models.CalendarEvents, req *CalendarUpdateRequest) {
	if req.Title != nil {
		cal.Title = *req.Title
	}
	if req.Description != nil {
		cal.Description = *req.Description
	}
	if req.Emoji != nil {
		cal.Emoji = *req.Emoji
	}
	if req.StartAt != nil {
		cal.StartAt = parseDate(*req.StartAt)
	}
	if req.EndAt != nil {
		cal.EndAt = parseDate(*req.EndAt)
	}
	if req.Visibility != nil {
		cal.Visibility = *req.Visibility
	}
	if req.ImageURL != nil {
		cal.ImageURL = stringOrEmpty(req.ImageURL)
	}
	if req.Todos != nil {
		todos := make([]*models.Todos, len(*req.Todos))
		for i, t := range *req.Todos {
			todos[i] = &models.Todos{
				Content: t.Text,
				IsDone:  t.Completed,
			}
		}
		cal.Todos = todos
	}
}

// ---------------------- 헬퍼 함수 ----------------------

func formatDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02")
}

func parseDate(s string) time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return t
}

func stringOrEmpty(ptr *string) string {
	if ptr == nil {
		return ""
	}
	return *ptr
}
