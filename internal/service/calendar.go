package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/mapper"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_user_server/internal/tx"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type CalendarServiceInterface interface {
	CreateCalendarEvent(ctx context.Context, cal *models.CalendarEvent) (*models.CalendarEvent, error)
	GetMyCalendarEvents(
		ctx context.Context,
		userID uuid.UUID,
		year, month int) ([]*models.CalendarEvent, error)
	GetEventDetailWithTodosByID(
		ctx context.Context,
		userID uuid.UUID,
		eventID uuid.UUID) (*models.CalendarEvent, error)
	GetMyCalendarEventsByDate(
		ctx context.Context,
		userID uuid.UUID,
		date time.Time) ([]*model.Calendar, error)
	DeleteCalendarEvent(
		ctx context.Context,
		UserID uuid.UUID,
		eventID uuid.UUID) error
	UpdateCalendarEvent(
		ctx context.Context,
		userID uuid.UUID,
		eventID uuid.UUID,
		input model.UpdateCalendarInput,
	) (*models.CalendarEvent, error)
}

type CalendarService struct {
	DB                 *gorm.DB
	ProfilesRepo       *repository.ProfileRepository
	CalendarEventsRepo *repository.CalendarEventsRepository
	// TodosRepo          *repository.TodosRepository
	// FollowsRepo        *repository.FollowsRepoitory
}

func NewCalendarService(
	db *gorm.DB,
	profilesRepo *repository.ProfileRepository,
	calendarRepo *repository.CalendarEventsRepository,
	// todoRepo *repository.TodosRepository,
) CalendarServiceInterface {
	return &CalendarService{
		DB:                 db,
		ProfilesRepo:       profilesRepo,
		CalendarEventsRepo: calendarRepo,
		// TodosRepo:          todoRepo,
	}
}

// // 내 캘린더 조회 (월별, Event만)
func (s *CalendarService) GetMyCalendarEvents(
	ctx context.Context,
	userID uuid.UUID,
	year, month int,
) ([]*models.CalendarEvent, error) {

	logger.Infof(
		"[GetMyCalendarEvents] userID=%s year=%d month=%d",
		userID, year, month,
	)

	start := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(0, 1, 0)

	return s.GetEventsWithoutTodos(
		ctx,
		userID,
		[]string{"public", "friends", "private"},
		start,
		end,
	)
}

func (s *CalendarService) GetEventDetailWithTodosByID(
	ctx context.Context,
	userID uuid.UUID,
	eventID uuid.UUID,
) (*models.CalendarEvent, error) {

	logger.Infof("[GetEventDetailWithTodosByID] eventID=%s userID=%s", eventID, userID)

	event, err := s.CalendarEventsRepo.GetEventWithTodosByID(ctx, eventID)
	if err != nil {
		return nil, fmt.Errorf("event not found or query failed: %w", err)
	}

	if event.UserID != userID {
		return nil, errors.New("forbidden: not your calendar event")
	}

	return event, nil
}

// // ----------------------------
// // Handler용 고수준 함수 (일별/Plan 전용)
// // ----------------------------

// // 내 일일 계획 조회 (일별, Event + Todo 포함, PlanHandler에서 호출)
func (s *CalendarService) GetMyCalendarEventsByDate(ctx context.Context, userID uuid.UUID, date time.Time) ([]*model.Calendar, error) {
	logger.Infof("[GetMyCalendarEventsByDate] UserID=%s, date=%s", userID, date.Format("2006-01-02"))

	// 조회 범위: 해당 일 00:00:00 부터 다음 날 00:00:00 까지
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)

	// Event와 Todo를 모두 포함하여 DB에서 조회 (캐시 미사용)
	calendars, err := s.CalendarEventsRepo.FindCalendarsWithTodos(ctx, userID, []string{"public", "friends", "private"}, startDate, endDate)
	if err != nil {
		logger.Errorf("[GetMyCalendarEventsByDate] FindCalendarsWithTodos failed: %v", err)
		return nil, errors.New("failed to get calendar events")
	}
	for _, v := range calendars {
		for _, p := range v.Todos {
			logger.Debugf("P : %v", p)
		}
	}
	// GraphQL 모델 변환
	return mapper.ToCalendarGraphQLList(calendars), nil
}

// GetEventsWithoutTodos: 월별 캘린더 뷰를 위해 Todo가 없는 Event만 조회하고 캐시 사용
func (s *CalendarService) GetEventsWithoutTodos(ctx context.Context, UserID uuid.UUID, visibilityLevels []string, startDate, endDate time.Time) ([]*models.CalendarEvent, error) {
	logger.Infof("[GetEventsWithoutTodos] user=%s, start=%s, end=%s", UserID, startDate, endDate)
	var allCalendars []*models.CalendarEvent

	remainingVis := make([]string, 0)

	for _, vis := range visibilityLevels {
		remainingVis = append(remainingVis, vis)
	}

	if len(remainingVis) > 0 {
		// 💡 Repository에서 Todo 없이 Event만 조회하는 메서드 사용 가정
		dbCalendars, err := s.CalendarEventsRepo.FindEventsWithoutTodosByVisibility(ctx, UserID, remainingVis, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("db error: %w", err)
		}

		for _, vis := range remainingVis {
			filtered := make([]*models.CalendarEvent, 0)
			for _, c := range dbCalendars {
				if c.Visibility == vis {
					filtered = append(filtered, c)
				}
			}
			allCalendars = append(allCalendars, filtered...)
		}
	}

	logger.Infof("[GetEventsWithoutTodos] retrieved %d events", len(allCalendars))
	return allCalendars, nil
}

// // ----------------------------
// // 기본 CRUD (캐시 무효화 포함)
// // ----------------------------

func (s *CalendarService) CreateCalendarEvent(
	ctx context.Context,
	cal *models.CalendarEvent,
) (*models.CalendarEvent, error) {

	txDB, newCtx, err := tx.BeginTx(ctx, s.DB)
	if err != nil {
		logger.Errorf("[CreateCalendar] failed to start transaction: %v", err)
		return nil, errors.New("failed to start transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[CreateCalendar] panic occurred, rollback: %v", r)
			txDB.Rollback()
			panic(r)
		}
	}()

	ctx = newCtx

	logger.Infof(
		"[CreateCalendar] start user=%s title=%s",
		cal.UserID,
		cal.Title,
	)

	// CalendarEvent 생성
	created, err := s.CalendarEventsRepo.CreateCalendarEvent(ctx, cal)
	if err != nil {
		txDB.Rollback()
		return nil, err
	}

	// Commit
	if err := txDB.Commit().Error; err != nil {
		logger.Errorf("[CreateCalendar] commit failed: %v", err)
		return nil, err
	}

	logger.Infof(
		"[CreateCalendar] success calendar_event_id=%s",
		created.ID,
	)

	return created, nil
}

func (s *CalendarService) UpdateCalendarEvent(
	ctx context.Context,
	userID uuid.UUID,
	eventID uuid.UUID,
	input model.UpdateCalendarInput,
) (*models.CalendarEvent, error) {

	event, err := s.CalendarEventsRepo.GetEventWithTodosByID(ctx, eventID)
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

	if err := s.CalendarEventsRepo.Update(ctx, event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *CalendarService) DeleteCalendarEvent(ctx context.Context, UserID uuid.UUID, eventID uuid.UUID) error {
	cal, err := s.CalendarEventsRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if cal == nil || cal.UserID != UserID {
		return fmt.Errorf("unauthorized or not found")
	}

	if err := s.CalendarEventsRepo.DeleteCalendarEvent(ctx, eventID); err != nil {
		return err
	}

	return nil
}

// // ----------------------------
// // Utility
// // ----------------------------

func (s *CalendarService) GenerateMonthData(startDate time.Time) [][]int {
	monthData := make([][]int, 6)
	for i := range monthData {
		monthData[i] = make([]int, 7)
	}

	firstWeekday := int(startDate.Weekday())
	daysInMonth := time.Date(startDate.Year(), startDate.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

	day := 1
weekLoop:
	for i := 0; i < len(monthData); i++ {
		for j := 0; j < 7; j++ {
			if i == 0 && j < firstWeekday {
				continue
			}
			if day > daysInMonth {
				break weekLoop
			}
			monthData[i][j] = day
			day++
		}
	}

	return monthData
}
