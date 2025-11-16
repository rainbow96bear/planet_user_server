package dto

import (
	"time"

	"github.com/rainbow96bear/planet_utils/model"
)

// ---------------------- DTO 구조 ----------------------

type CalendarInfo struct {
	EventID     uint64     `json:"eventId"`
	UserUUID    string     `json:"userUuid"`
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
	Title       string     `json:"title" form:"title" binding:"required,max=255"`
	Description string     `json:"description" form:"description" binding:"max=65535"`
	Emoji       string     `json:"emoji" form:"emoji" binding:"omitempty,max=10"`
	StartAt     string     `json:"startAt" form:"startAt" binding:"required,datetime=2006-01-02"`
	EndAt       string     `json:"endAt" form:"endAt" binding:"required,datetime=2006-01-02,gtefield=StartAt"`
	Visibility  string     `json:"visibility" form:"visibility" binding:"required,oneof=public friends private"`
	Todos       []TodoItem `json:"todos" form:"todos" binding:"omitempty,dive"`
	ImageURL    *string    `json:"imageUrl" form:"imageUrl"`
}

type CalendarUpdateRequest struct {
	Title       *string     `json:"title,omitempty" form:"title" binding:"omitempty,max=255"`
	Description *string     `json:"description,omitempty" form:"description" binding:"omitempty,max=65535"`
	Emoji       *string     `json:"emoji,omitempty" form:"emoji" binding:"omitempty,max=10"`
	StartAt     *string     `json:"startAt,omitempty" form:"startAt" binding:"omitempty,datetime=2006-01-02"`
	EndAt       *string     `json:"endAt,omitempty" form:"endAt" binding:"omitempty,datetime=2006-01-02"`
	Visibility  *string     `json:"visibility,omitempty" form:"visibility" binding:"omitempty,oneof=public friends private"`
	Todos       *[]TodoItem `json:"todos,omitempty" form:"todos" binding:"omitempty,dive"`
	ImageURL    *string     `json:"imageUrl,omitempty" form:"imageUrl"`
}

type TodoItem struct {
	EventID   *uint64 `json:"eventId,omitempty" form:"eventId"`
	Text      string  `json:"text" form:"text" binding:"required,max=255"`
	Completed bool    `json:"completed" form:"completed"`
}

// ---------------------- 변환 함수 ----------------------

func ToCalendarInfo(cal *model.Calendar) *CalendarInfo {
	if cal == nil {
		return nil
	}

	todos := make([]TodoItem, len(cal.Todos))
	for i, t := range cal.Todos {
		todos[i] = TodoItem{
			EventID:   &t.EventID,
			Text:      t.Content,
			Completed: t.Done,
		}
	}

	return &CalendarInfo{
		EventID:     cal.EventID,
		UserUUID:    cal.UserUUID,
		Title:       cal.Title,
		Description: cal.Description,
		Emoji:       cal.Emoji,
		StartAt:     formatDate(cal.StartAt),
		EndAt:       formatDate(cal.EndAt),
		Visibility:  cal.Visibility,
		ImageURL:    cal.ImageURL,
		Todos:       todos,
		CreatedAt:   cal.CreatedAt,
		UpdatedAt:   cal.UpdatedAt,
	}
}

func ToCalendarInfoList(models []*model.Calendar) []*CalendarInfo {
	result := make([]*CalendarInfo, 0, len(models))
	for _, m := range models {
		if info := ToCalendarInfo(m); info != nil {
			result = append(result, info)
		}
	}
	return result
}

// Create 요청을 GORM 모델로 변환
func ToCalendarModelFromCreate(req *CalendarCreateRequest, userUUID string) *model.Calendar {
	cal := &model.Calendar{
		UserUUID:    userUUID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     parseDate(req.StartAt),
		EndAt:       parseDate(req.EndAt),
		Visibility:  req.Visibility,
		ImageURL:    req.ImageURL,
		Todos:       make([]model.Todo, len(req.Todos)),
	}

	for i, t := range req.Todos {
		cal.Todos[i] = model.Todo{
			Content: t.Text,
			Done:    t.Completed,
		}
	}

	return cal
}

// Update 요청을 기존 모델에 적용 (GORM용)
func UpdateCalendarModelFromRequest(cal *model.Calendar, req *CalendarUpdateRequest) {
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
		cal.ImageURL = req.ImageURL
	}

	// Todos 업데이트
	if req.Todos != nil {
		var todos []model.Todo
		for _, t := range *req.Todos {
			todo := model.Todo{
				Content: t.Text,
				Done:    t.Completed,
			}
			todos = append(todos, todo)
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
