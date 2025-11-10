package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_utils/model"
)

type CalendarService struct {
	CalendarRepo *repository.CalendarRepository
	UsersRepo    *repository.UsersRepository
}

// -------------------- 유틸 --------------------
func validateCalendarFields(c *dto.CalendarInfo) error {
	if c == nil {
		return fmt.Errorf("invalid calendar data")
	}
	if c.Title == "" {
		return fmt.Errorf("title is required")
	}
	if c.StartAt == "" || c.EndAt == "" {
		return fmt.Errorf("start date and end date are required")
	}
	if c.Visibility != "public" && c.Visibility != "friends" && c.Visibility != "private" {
		return fmt.Errorf("invalid visibility value")
	}
	return nil
}

// DTO → Model 변환
func toModel(c *dto.CalendarInfo) *model.Calendar {
	startAt, _ := time.Parse("2006-01-02", c.StartAt)
	endAt, _ := time.Parse("2006-01-02", c.EndAt)

	todos := make([]model.Todo, len(c.Todos))
	for i, t := range c.Todos {
		todos[i] = model.Todo{
			Content: t.Text,
			Done:    t.Completed,
		}
	}

	return &model.Calendar{
		ID:          c.EventID,
		UserUUID:    c.UserUUID,
		Title:       c.Title,
		Description: c.Description,
		Emoji:       c.Emoji,
		StartAt:     startAt,
		EndAt:       endAt,
		Visibility:  c.Visibility,
		ImageURL:    &c.ImageURL,
		Todos:       todos,
	}
}

// Model → DTO 변환
func toDTO(m *model.Calendar) *dto.CalendarInfo {
	todos := make([]dto.TodoItem, len(m.Todos))
	for i, t := range m.Todos {
		todos[i] = dto.TodoItem{
			Text:      t.Content,
			Completed: t.Done,
		}
	}

	imageURL := ""
	if m.ImageURL != nil {
		imageURL = *m.ImageURL
	}

	return &dto.CalendarInfo{
		EventID:     m.ID,
		UserUUID:    m.UserUUID,
		Title:       m.Title,
		Description: m.Description,
		Emoji:       m.Emoji,
		StartAt:     m.StartAt.Format("2006-01-02"),
		EndAt:       m.EndAt.Format("2006-01-02"),
		Visibility:  m.Visibility,
		ImageURL:    imageURL,
		Todos:       todos,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func convertModelListToDTO(models []*model.Calendar) []*dto.CalendarInfo {
	result := make([]*dto.CalendarInfo, len(models))
	for i, m := range models {
		result[i] = toDTO(m)
	}
	return result
}

// -------------------- 조회 --------------------
func (s *CalendarService) GetUserCalendar(ctx context.Context, userUUID string) ([]*dto.CalendarInfo, error) {
	models, err := s.CalendarRepo.GetCalendarsByUserUuid(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return convertModelListToDTO(models), nil
}

func (s *CalendarService) GetCalendarByID(ctx context.Context, userUUID string, eventID uint64) (*dto.CalendarInfo, error) {
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, eventID, userUUID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, fmt.Errorf("unauthorized to access this event")
	}

	m, err := s.CalendarRepo.GetCalendarByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return toDTO(m), nil
}

// -------------------- 생성/수정/삭제 --------------------
func (s *CalendarService) CreateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if err := validateCalendarFields(calendar); err != nil {
		return err
	}

	model := toModel(calendar)
	return s.CalendarRepo.CreateCalendar(ctx, model)
}

func (s *CalendarService) UpdateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	if err := validateCalendarFields(calendar); err != nil {
		return err
	}

	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, calendar.EventID, calendar.UserUUID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to update this event")
	}

	model := toModel(calendar)
	return s.CalendarRepo.UpdateCalendar(ctx, model)
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, userUUID string, eventID uint64) error {
	isOwner, err := s.CalendarRepo.IsOwnerOfCalendar(ctx, eventID, userUUID)
	if err != nil {
		return err
	}
	if !isOwner {
		return fmt.Errorf("unauthorized to delete this event")
	}
	return s.CalendarRepo.DeleteCalendar(ctx, eventID)
}

// -------------------- 공개 캘린더 조회 --------------------
func (s *CalendarService) GetPublicCalendarByNickname(ctx context.Context, nickname string) ([]*dto.CalendarInfo, error) {
	userUUID, err := s.UsersRepo.GetUserUuidByNickname(ctx, nickname)
	if err != nil {
		return nil, err
	}

	models, err := s.CalendarRepo.GetPublicCalendarsByUserUuid(ctx, userUUID)
	if err != nil {
		return nil, err
	}
	return convertModelListToDTO(models), nil
}

func (s *CalendarService) GetAllPublicCalendars(ctx context.Context) ([]*dto.CalendarInfo, error) {
	models, err := s.CalendarRepo.GetAllPublicCalendars(ctx)
	if err != nil {
		return nil, err
	}
	return convertModelListToDTO(models), nil
}

func (s *CalendarService) GetCalendarByNicknameAndVisibility(
	ctx context.Context,
	nickname string,
	visibilityLevels []string,
) ([]*dto.CalendarInfo, error) {
	models, err := s.CalendarRepo.FindByNicknameAndVisibility(ctx, nickname, visibilityLevels)
	if err != nil {
		return nil, err
	}
	return convertModelListToDTO(models), nil
}
