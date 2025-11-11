package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_utils/model"
)

type CalendarService struct {
	CalendarRepo *repository.CalendarRepository
	UsersRepo    *repository.UsersRepository
}

func (s *CalendarService) GetUserCalendars(
	ctx context.Context,
	userUUID string,
	visibilityLevels []string,
	startDate, endDate time.Time,
) ([]*model.Calendar, error) {
	var allCalendars []*model.Calendar

	// 먼저 캐시에서 visibility별로 가져오기
	remainingVis := make([]string, 0)
	for _, vis := range visibilityLevels {
		if cached, ok := GetCalendarCache(userUUID, startDate.Year(), int(startDate.Month()), vis); ok {
			allCalendars = append(allCalendars, cached...)
		} else {
			remainingVis = append(remainingVis, vis)
		}
	}

	// 남은 visibility가 없으면 바로 반환
	if len(remainingVis) == 0 {
		return allCalendars, nil
	}

	// ✅ DB 조회: visibility를 배열로 전달
	calendars, err := s.CalendarRepo.FindCalendarsByVisibility(ctx, userUUID, remainingVis, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	// 조회한 데이터를 visibility별로 분리하여 캐시 저장
	for _, vis := range remainingVis {
		filtered := make([]*model.Calendar, 0)
		for _, c := range calendars {
			if c.Visibility == vis {
				filtered = append(filtered, c)
			}
		}
		SetCalendarCache(userUUID, startDate.Year(), int(startDate.Month()), vis, filtered)
		allCalendars = append(allCalendars, filtered...)
	}

	return allCalendars, nil
}

func (s *CalendarService) CreateCalendar(ctx context.Context, cal *model.Calendar) error {
	// Transaction: Calendar + Todos
	return s.CalendarRepo.CreateCalendarWithTodos(ctx, cal)
}

func (s *CalendarService) DeleteCalendar(ctx context.Context, userUUID string, eventID uint64) error {
	// 1. DB에서 삭제
	cal, err := s.CalendarRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if cal.UserUUID != userUUID {
		return fmt.Errorf("unauthorized")
	}

	if err := s.CalendarRepo.DeleteCalendar(ctx, eventID); err != nil {
		return err
	}

	// 2. 캐시 삭제
	yearStart := cal.StartAt.Year()
	monthStart := int(cal.StartAt.Month())
	DeleteCalendarCache(userUUID, yearStart, monthStart, cal.Visibility)

	return nil
}

func (s *CalendarService) ClearCache(userUUID string, year, month int) {
	for _, vis := range []string{"public", "friends", "private"} {
		DeleteCalendarCache(userUUID, year, month, vis)
	}
}

func (s *CalendarService) GenerateMonthData(startDate time.Time) [][]int {
	monthData := make([][]int, 6)
	for i := range monthData {
		monthData[i] = make([]int, 7)
	}

	firstWeekday := int(startDate.Weekday()) // 0 = Sunday
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

func (s *CalendarService) CalculateCompletionData(calendars []*model.Calendar) map[int]int {
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
			if t.Done {
				doneCount++
			}
		}
		completion[day] = doneCount * 100 / totalTodos
	}
	return completion
}
