package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"gorm.io/gorm"
)

type CalendarService struct {
	db           *gorm.DB
	CalendarRepo *repository.CalendarEventsRepository
	FollowRepo   *repository.FollowsRepository
}

func NewCalendarService(
	db *gorm.DB,
	calendarRepo *repository.CalendarEventsRepository,
	followRepo *repository.FollowsRepository,
) *CalendarService {
	return &CalendarService{
		db:           db,
		CalendarRepo: calendarRepo,
		FollowRepo:   followRepo,
	}
}

//
// =======================
// Create
// =======================
//

func (s *CalendarService) Create(
	ctx context.Context,
	cal *models.CalendarEvent,
) (*models.CalendarEvent, error) {

	txDB, newCtx, err := tx.BeginTx(ctx, s.db)
	if err != nil {
		return nil, errors.New("failed to start transaction")
	}
	ctx = newCtx

	defer func() {
		if r := recover(); r != nil {
			txDB.Rollback()
			panic(r)
		}
	}()

	created, err := s.CalendarRepo.CreateCalendarEvent(ctx, cal)
	if err != nil {
		txDB.Rollback()
		return nil, err
	}

	if err := txDB.Commit().Error; err != nil {
		return nil, err
	}

	return created, nil
}

//
// =======================
// Update
// =======================
//

func (s *CalendarService) Update(
	ctx context.Context,
	userID uuid.UUID,
	eventID uuid.UUID,
	input model.UpdateCalendarInput,
) (*models.CalendarEvent, error) {

	event, err := s.CalendarRepo.GetEventWithTodosByID(ctx, eventID)
	if err != nil {
		return nil, errors.New("event not found")
	}

	if event.UserID != userID {
		return nil, errors.New("forbidden")
	}

	req := dto.CalendarUpdateRequest{
		Title:       input.Title,
		Emoji:       input.Emoji,
		Description: input.Description,
		StartAt:     input.StartAt,
		EndAt:       input.EndAt,
		Visibility:  (*string)(input.Visibility),
	}

	if input.Todos != nil {
		req.Todos = make([]dto.TodoUpdateRequest, 0, len(input.Todos))
		for _, t := range input.Todos {
			req.Todos = append(req.Todos, dto.TodoUpdateRequest{
				Content: t.Content,
				IsDone:  t.IsDone,
			})
		}
	}

	dto.UpdateCalendarModelFromRequest(event, &req)

	if err := s.CalendarRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

//
// =======================
// Delete
// =======================
//

func (s *CalendarService) Delete(
	ctx context.Context,
	userID uuid.UUID,
	eventID uuid.UUID,
) error {

	event, err := s.CalendarRepo.FindByID(ctx, eventID)
	if err != nil || event == nil {
		return errors.New("event not found")
	}

	if event.UserID != userID {
		return errors.New("forbidden")
	}

	return s.CalendarRepo.DeleteCalendarEvent(ctx, eventID)
}

//
// =======================
// Query
// =======================
//

func (s *CalendarService) GetMyEventsByPeriod(
	ctx context.Context,
	userID uuid.UUID,
	from, to time.Time,
) ([]*models.CalendarEvent, error) {

	return s.CalendarRepo.FindEventsWithoutTodosByVisibility(
		ctx,
		userID,
		[]string{"public", "friends", "private"},
		from,
		to,
	)
}

func (s *CalendarService) GetUserEventsByPeriod(
	ctx context.Context,
	viewerID *uuid.UUID,
	targetUserID uuid.UUID,
	from, to time.Time,
) ([]*models.CalendarEvent, error) {

	// 기본: 비로그인 or 타인 → public only
	visibilities := []string{"public"}

	// 로그인 상태일 때만 친구 여부 체크
	if viewerID != nil {
		isFriend, err := s.FollowRepo.IsFollowing(
			ctx,
			*viewerID,
			targetUserID,
		)
		if err != nil {
			return nil, err
		}

		if isFriend {
			visibilities = append(visibilities, "friends")
		}
	}

	return s.CalendarRepo.FindEventsWithoutTodosByVisibility(
		ctx,
		targetUserID,
		visibilities,
		from,
		to,
	)
}

func (s *CalendarService) GetDetailWithTodos(
	ctx context.Context,
	userID uuid.UUID,
	eventID uuid.UUID,
) (*models.CalendarEvent, error) {

	event, err := s.CalendarRepo.GetEventWithTodosByID(ctx, eventID)
	if err != nil {
		return nil, err
	}

	if event.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return event, nil
}

func (s *CalendarService) GetMyEventsByDateWithTodos(
	ctx context.Context,
	userID uuid.UUID,
	date time.Time,
) ([]*models.CalendarEvent, error) {

	start := time.Date(
		date.Year(),
		date.Month(),
		date.Day(),
		0, 0, 0, 0,
		time.UTC,
	)
	end := start.AddDate(0, 0, 1)

	return s.CalendarRepo.FindCalendarsWithTodos(
		ctx,
		userID,
		[]string{"public", "friends", "private"},
		start,
		end,
	)
}
