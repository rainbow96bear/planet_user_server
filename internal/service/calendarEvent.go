package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_utils/models"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type CalendarService struct {
	CalendarEventsRepo *repository.CalendarEventsRepository
	ProfilesRepo       *repository.ProfilesRepository
	FollowsRepo        *repository.FollowsRepository
}

// ----------------------------
// Handler용 고수준 함수
// ----------------------------

// 다른 사람 캘린더 조회 (visibility 판단 포함)
func (s *CalendarService) GetUserCalendarData(ctx context.Context, nickname string, authID uuid.UUID, year, month int) (map[string]interface{}, error) {
	logger.Infof("[GetUserCalendarData] nickname=%s, authUUID=%s, year=%d month=%d", nickname, authID, year, month)

	// 사용자 UUID 조회 (Repository)
	UserID, err := s.ProfilesRepo.GetUserIDByNickname(ctx, nickname)
	if err != nil {
		logger.Errorf("[GetUserCalendarData] failed to get user UUID: %v", err)
		return nil, err
	}

	// visibility 결정
	visibility := []string{"public"}
	if authID != uuid.Nil && authID != UserID {
		isFollow, _ := s.FollowsRepo.IsFollow(ctx, authID, UserID)
		if isFollow {
			visibility = append(visibility, "friends")
		}
	} else if authID == UserID {
		visibility = append(visibility, "friends", "private")
	}

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	calendars, err := s.GetUserCalendars(ctx, UserID, visibility, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events":         ToDTOList(calendars),
		"monthData":      s.GenerateMonthData(startDate),
		"completionData": s.CalculateCompletionData(calendars),
	}, nil
}

// 내 캘린더 조회
func (s *CalendarService) GetMyCalendarData(ctx context.Context, UserID uuid.UUID, year, month int) (map[string]interface{}, error) {
	logger.Infof("[GetMyCalendarData] UserID=%s, year=%d month=%d", UserID, year, month)

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	calendars, err := s.GetUserCalendars(ctx, UserID, []string{"public", "friends", "private"}, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events":         ToDTOList(calendars),
		"monthData":      s.GenerateMonthData(startDate),
		"completionData": s.CalculateCompletionData(calendars),
	}, nil
}

// ----------------------------
// 기본 CRUD
// ----------------------------

func (s *CalendarService) GetUserCalendars(ctx context.Context, UserID uuid.UUID, visibilityLevels []string, startDate, endDate time.Time) ([]*models.CalendarEvents, error) {
	logger.Infof("[GetUserCalendars] user=%s, start=%s, end=%s", UserID, startDate, endDate)
	var allCalendars []*models.CalendarEvents

	remainingVis := make([]string, 0)

	// 캐시 조회
	for _, vis := range visibilityLevels {
		if cached, ok := GetCalendarCache(UserID, startDate.Year(), int(startDate.Month()), vis); ok {
			allCalendars = append(allCalendars, cached...)
		} else {
			remainingVis = append(remainingVis, vis)
		}
	}

	if len(remainingVis) > 0 {
		dbCalendars, err := s.CalendarEventsRepo.FindCalendarsByVisibility(ctx, UserID, remainingVis, startDate, endDate)
		if err != nil {
			return nil, fmt.Errorf("db error: %w", err)
		}

		for _, vis := range remainingVis {
			filtered := make([]*models.CalendarEvents, 0)
			for _, c := range dbCalendars {
				if c.Visibility == vis {
					filtered = append(filtered, c)
				}
			}
			SetCalendarCache(UserID, startDate.Year(), int(startDate.Month()), vis, filtered)
			allCalendars = append(allCalendars, filtered...)
		}
	}

	logger.Infof("[GetUserCalendars] retrieved %d calendars", len(allCalendars))
	return allCalendars, nil
}

func (s *CalendarService) CreateCalendar(ctx context.Context, cal *models.CalendarEvents) error {
	logger.Infof("[CreateCalendar] user=%s title=%s", cal.UserID, cal.Title)
	if err := s.CalendarEventsRepo.CreateCalendarEvent(ctx, cal); err != nil {
		logger.Errorf("[CreateCalendar] failed: %v", err)
		return err
	}
	ClearCache(cal.UserID, cal.StartAt.Year(), int(cal.StartAt.Month()))
	logger.Infof("[CreateCalendar] successfully created calendar event: %s", cal.ID)
	return nil
}

func (s *CalendarService) UpdateCalendar(ctx context.Context, UserID uuid.UUID, eventID uuid.UUID, req *dto.CalendarUpdateRequest) error {
	cal, err := s.CalendarEventsRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if cal.UserID != UserID {
		return fmt.Errorf("unauthorized")
	}

	dto.UpdateCalendarModelFromRequest(cal, req)

	if err := s.CalendarEventsRepo.UpdateCalendar(ctx, cal); err != nil {
		return err
	}

	ClearCache(cal.UserID, cal.StartAt.Year(), int(cal.StartAt.Month()))
	return nil
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, UserID uuid.UUID, eventID uuid.UUID) error {
	cal, err := s.CalendarEventsRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if cal == nil || cal.UserID != UserID {
		return fmt.Errorf("unauthorized or not found")
	}

	if err := s.CalendarEventsRepo.DeleteCalendar(ctx, eventID); err != nil {
		return err
	}

	DeleteCalendarCache(UserID, cal.StartAt.Year(), int(cal.StartAt.Month()), cal.Visibility)
	return nil
}

// ----------------------------
// Utility
// ----------------------------

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

func (s *CalendarService) CalculateCompletionData(calendars []*models.CalendarEvents) map[int]int {
	completion := make(map[int]int)
	for _, cal := range calendars {
		day := cal.StartAt.Day()
		totalTodos := len(cal.Todos)
		if totalTodos == 0 {
			completion[day] = 100
			continue
		}
		doneCount := 0
		for _, t := range cal.Todos {
			if t.IsDone {
				doneCount++
			}
		}
		completion[day] = doneCount * 100 / totalTodos
	}
	return completion
}

// ----------------------------
// DTO 변환 헬퍼
// ----------------------------
func ToDTOList(calendars []*models.CalendarEvents) []*dto.CalendarInfo {
	result := make([]*dto.CalendarInfo, 0, len(calendars))
	for _, cal := range calendars {
		result = append(result, dto.ToCalendarInfo(cal))
	}
	return result
}

// 전체 캐시 초기화
func ClearCache(UserID uuid.UUID, year, month int) {
	for _, vis := range []string{"public", "friends", "private"} {
		DeleteCalendarCache(UserID, year, month, vis)
	}
}
