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
	TodosRepo          *repository.TodosRepository
	ProfilesRepo       *repository.ProfilesRepository
	FollowsRepo        *repository.FollowsRepository
}

// ----------------------------
// Handler용 고수준 함수 (월별/Event 전용)
// ----------------------------

// 다른 사람 캘린더 조회 (월별, Event만)
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

	// 💡 Todo가 없는 Event만 조회 (캐시 활용)
	calendars, err := s.GetEventsWithoutTodos(ctx, UserID, visibility, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events":    ToDTOList(calendars),
		"monthData": s.GenerateMonthData(startDate),
		// "completionData": 월별 조회에서는 Todo가 없으므로 반환하지 않음
	}, nil
}

// 내 캘린더 조회 (월별, Event만)
func (s *CalendarService) GetMyCalendarData(ctx context.Context, UserID uuid.UUID, year, month int) (map[string]interface{}, error) {
	logger.Infof("[GetMyCalendarData] UserID=%s, year=%d month=%d", UserID, year, month)

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0)

	// 💡 Todo가 없는 Event만 조회 (캐시 활용)
	calendars, err := s.GetEventsWithoutTodos(ctx, UserID, []string{"public", "friends", "private"}, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"events":    ToDTOList(calendars),
		"monthData": s.GenerateMonthData(startDate),
		// "completionData": 월별 조회에서는 Todo가 없으므로 반환하지 않음
	}, nil
}

func (s *CalendarService) GetEventDetailWithTodosByID(ctx context.Context, eventID uuid.UUID) (*dto.CalendarInfo, error) {
	// 💡 UserID 매개변수 제거: 권한 확인을 수행하지 않으므로 필요하지 않습니다.
	logger.Infof("[GetEventDetailWithTodosByID] EventID=%d", eventID)

	// 1. Repository 호출: eventID로 이벤트와 Todo를 함께 조회합니다.
	event, err := s.CalendarEventsRepo.FindEventWithTodosByID(ctx, eventID)

	if err != nil {
		// DB 조회 실패 (예: 해당 eventID의 레코드가 없는 경우)
		// DTO 반환 전에 에러를 처리하여 상위 계층에 전달합니다.
		return nil, fmt.Errorf("event not found or query failed for ID %d: %w", eventID, err)
	}

	// 2. DTO로 변환 및 반환
	// event 모델에 이미 Todos가 로드되어 있다고 가정하고 DTO로 변환합니다.
	eventDTO := dto.ToCalendarInfo(event)

	// 💡 참고: 권한(UserID/Visibility) 확인 로직은 이 함수에서 완전히 제거되었습니다.

	return eventDTO, nil
}

// ----------------------------
// Handler용 고수준 함수 (일별/Plan 전용)
// ----------------------------

// 내 일일 계획 조회 (일별, Event + Todo 포함, PlanHandler에서 호출)
func (s *CalendarService) GetMyCalendarDailyData(ctx context.Context, UserID uuid.UUID, date time.Time) (map[string]interface{}, error) {
	logger.Infof("[GetMyCalendarDailyData] UserID=%s, date=%s", UserID, date.Format("2006-01-02"))

	// 조회 범위: 해당 일 00:00:00 부터 다음 날 00:00:00 까지
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)

	// 💡 Event와 Todo를 모두 포함하여 DB에서 조회 (캐시 미사용)
	calendars, err := s.CalendarEventsRepo.FindCalendarsWithTodos(ctx, UserID, []string{"public", "friends", "private"}, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// CalculateCompletionData를 사용하여 일별 달성률 데이터를 포함할 수 있습니다.
	completionData := s.CalculateCompletionData(calendars)

	return map[string]interface{}{
		"dailyPlans":     ToDTOList(calendars),
		"completionData": completionData,
	}, nil
}

// 다른 사람 일일 계획 조회 (일별, Event + Todo 포함, PlanHandler에서 호출)
func (s *CalendarService) GetUserCalendarDailyData(ctx context.Context, nickname string, authID uuid.UUID, date time.Time) (map[string]interface{}, error) {
	logger.Infof("[GetUserCalendarDailyData] nickname=%s, authUUID=%s, date=%s", nickname, authID, date.Format("2006-01-02"))

	// 1. 사용자 UUID 조회 및 Visibility 결정
	UserID, err := s.ProfilesRepo.GetUserIDByNickname(ctx, nickname)
	if err != nil {
		logger.Errorf("[GetUserCalendarDailyData] failed to get user UUID: %v", err)
		return nil, err
	}

	visibility := []string{"public"}
	if authID != uuid.Nil && authID != UserID {
		isFollow, _ := s.FollowsRepo.IsFollow(ctx, authID, UserID)
		if isFollow {
			visibility = append(visibility, "friends")
		}
	} else if authID == UserID {
		visibility = append(visibility, "friends", "private")
	}

	// 2. 조회 범위 설정
	startDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 0, 1)

	// 3. Event와 Todo를 모두 포함하여 DB에서 조회
	calendars, err := s.CalendarEventsRepo.FindCalendarsWithTodos(ctx, UserID, visibility, startDate, endDate)
	if err != nil {
		return nil, err
	}

	completionData := s.CalculateCompletionData(calendars)

	return map[string]interface{}{
		"dailyPlans":     ToDTOList(calendars),
		"completionData": completionData,
	}, nil
}

// ----------------------------
// Todo 상태 업데이트 (새로 추가)
// ----------------------------

func (s *CalendarService) UpdateTodoStatus(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, isDone bool) error {
	logger.Infof("[UpdateTodoStatus] UserID=%s, TodoID=%s, IsDone=%t", userID, todoID, isDone)

	// 1. Todo 상태 업데이트 및 소유권 확인
	err := s.TodosRepo.UpdateTodoStatus(ctx, todoID, isDone)
	if err != nil {
		return err
	}
	logger.Infof("[UpdateTodoStatus] Finished successfully: TodoID=%s", todoID)
	return nil
}

// ----------------------------
// 기본 조회 함수 (캐시 활용)
// ----------------------------

// GetEventsWithoutTodos: 월별 캘린더 뷰를 위해 Todo가 없는 Event만 조회하고 캐시 사용
func (s *CalendarService) GetEventsWithoutTodos(ctx context.Context, UserID uuid.UUID, visibilityLevels []string, startDate, endDate time.Time) ([]*models.CalendarEvents, error) {
	logger.Infof("[GetEventsWithoutTodos] user=%s, start=%s, end=%s", UserID, startDate, endDate)
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
		// 💡 Repository에서 Todo 없이 Event만 조회하는 메서드 사용 가정
		dbCalendars, err := s.CalendarEventsRepo.FindEventsWithoutTodosByVisibility(ctx, UserID, remainingVis, startDate, endDate)
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
			// 💡 Todo가 없는 Event만 캐시에 저장
			SetCalendarCache(UserID, startDate.Year(), int(startDate.Month()), vis, filtered)
			allCalendars = append(allCalendars, filtered...)
		}
	}

	logger.Infof("[GetEventsWithoutTodos] retrieved %d events", len(allCalendars))
	return allCalendars, nil
}

// ----------------------------
// 기본 CRUD (캐시 무효화 포함)
// ----------------------------

func (s *CalendarService) CreateCalendarEvent(ctx context.Context, cal *models.CalendarEvents) error {
	logger.Infof("[CreateCalendar] user=%s title=%s", cal.UserID, cal.Title)
	if err := s.CalendarEventsRepo.CreateCalendarEvent(ctx, cal); err != nil {
		logger.Errorf("[CreateCalendar] failed: %v", err)
		return err
	}
	// Event 생성 시 해당 월의 모든 가시성 캐시를 삭제
	ClearCache(cal.UserID, cal.StartAt.Year(), int(cal.StartAt.Month()))
	logger.Infof("[CreateCalendar] successfully created calendar event: %s", cal.ID)
	return nil
}

func (s *CalendarService) UpdateCalendarEvent(ctx context.Context, UserID uuid.UUID, eventID uuid.UUID, req *dto.CalendarUpdateRequest) error {
	cal, err := s.CalendarEventsRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	if cal.UserID != UserID {
		return fmt.Errorf("unauthorized")
	}

	dto.UpdateCalendarModelFromRequest(cal, req)

	if err := s.CalendarEventsRepo.UpdateCalendarEvent(ctx, cal); err != nil {
		return err
	}

	// Event 업데이트 시 해당 월의 모든 가시성 캐시를 삭제
	ClearCache(cal.UserID, cal.StartAt.Year(), int(cal.StartAt.Month()))
	return nil
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

	// Event 삭제 시 해당 월의 특정 가시성 캐시를 삭제
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

// CalculateCompletionData: 월별 조회에서는 사용되지 않지만, 일별 조회에서 사용됩니다.
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
