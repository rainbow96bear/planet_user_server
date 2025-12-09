package repository

import (
	"gorm.io/gorm"
)

type CalendarEventsRepository struct {
	DB *gorm.DB
}

// -------------------------
// 트랜잭션 시작
// -------------------------
// func (r *CalendarEventsRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
// 	logger.Infof("starting transaction for CalendarEventsRepository")
// 	tx := r.DB.WithContext(ctx).Begin()
// 	if tx.Error != nil {
// 		logger.Errorf("failed to start transaction: %v", tx.Error)
// 		return nil, tx.Error
// 	}
// 	logger.Infof("transaction started successfully")
// 	return tx, nil
// }

// // -------------------------
// // 캘린더 이벤트 생성 (Todos 포함)
// // -------------------------
// func (r *CalendarEventsRepository) CreateCalendarEvent(ctx context.Context, event *models.CalendarEvents) error {
// 	logger.Infof("Creating calendar event for user: %s", event.UserID)

// 	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// 1. Calendar Event 삽입
// 		if err := tx.Create(event).Error; err != nil {
// 			return fmt.Errorf("failed to insert calendar event: %w", err)
// 		}

// 		logger.Infof("Successfully created calendar event %s with %d todos", event.ID, len(event.Todos))
// 		return nil
// 	})
// }

// // -------------------------
// // 단일 조회 (Todos 포함)
// // -------------------------
// func (r *CalendarEventsRepository) FindByID(ctx context.Context, eventID uuid.UUID) (*models.CalendarEvents, error) {
// 	var event models.CalendarEvents
// 	if err := r.DB.WithContext(ctx).
// 		Preload("Todos"). // Todo도 함께 조회
// 		First(&event, "id = ?", eventID).Error; err != nil {
// 		if err == gorm.ErrRecordNotFound {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("failed to find calendar event: %w", err)
// 	}
// 	return &event, nil
// }

// // -------------------------
// // 캘린더 이벤트 삭제 (Todos 포함)
// // -------------------------
// func (r *CalendarEventsRepository) DeleteCalendarEvent(ctx context.Context, eventID uuid.UUID) error {
// 	logger.Infof("Deleting calendar event: %s", eventID)

// 	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// Todos 먼저 삭제 (Foreign Key 제약 조건)
// 		if err := tx.Where("event_id = ?", eventID).Delete(&models.Todos{}).Error; err != nil {
// 			return fmt.Errorf("failed to delete todos: %w", err)
// 		}
// 		// Event 삭제
// 		if err := tx.Where("id = ?", eventID).Delete(&models.CalendarEvents{}).Error; err != nil {
// 			return fmt.Errorf("failed to delete calendar event: %w", err)
// 		}
// 		logger.Infof("Deleted calendar event %s and its todos", eventID)
// 		return nil
// 	})
// }

// // -------------------------
// // 캘린더 이벤트 업데이트 (Todos 포함)
// // -------------------------
// func (r *CalendarEventsRepository) UpdateCalendarEvent(ctx context.Context, event *models.CalendarEvents) error {
// 	logger.Infof("[UpdateCalendar] eventID=%s", event.ID)

// 	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		// CalendarEvent 업데이트
// 		if err := tx.Save(event).Error; err != nil {
// 			logger.Errorf("[UpdateCalendar] failed to update event: %v", err)
// 			return fmt.Errorf("failed to update calendar event: %w", err)
// 		}

// 		// 기존 Todos 삭제 후 새로 삽입
// 		if err := tx.Where("event_id = ?", event.ID).Delete(&models.Todos{}).Error; err != nil {
// 			logger.Errorf("[UpdateCalendar] failed to delete old todos: %v", err)
// 			return fmt.Errorf("failed to delete old todos: %w", err)
// 		}
// 		for i := range event.Todos {
// 			// 업데이트 시에도 새 ID 할당 (혹은 기존 ID 재사용 로직 구현 필요하지만, 여기서는 단순화하여 새 삽입)
// 			event.Todos[i].ID = uuid.New()
// 			event.Todos[i].EventID = event.ID
// 		}
// 		if len(event.Todos) > 0 {
// 			if err := tx.Create(&event.Todos).Error; err != nil {
// 				logger.Errorf("[UpdateCalendar] failed to insert new todos: %v", err)
// 				return fmt.Errorf("failed to insert new todos: %w", err)
// 			}
// 		}

// 		logger.Infof("[UpdateCalendar] successfully updated eventID=%s with %d todos", event.ID, len(event.Todos))
// 		return nil
// 	})
// }

// // ------------------------------------------
// // 조회 함수 1: 월별 뷰 (Event만, 캐시 지원)
// // ------------------------------------------

// // FindEventsWithoutTodosByVisibility: 특정 기간 동안의 Event를 Todo 없이 조회합니다.
// // CalendarService의 GetEventsWithoutTodos에서 사용됩니다. (캐싱 목적)
// func (r *CalendarEventsRepository) FindEventsWithoutTodosByVisibility(
// 	ctx context.Context,
// 	UserID uuid.UUID,
// 	visibilities []string,
// 	startAt, endAt time.Time,
// ) ([]*models.CalendarEvents, error) {
// 	logger.Infof("Fetching events (without todos) for user=%s with visibilities=%v", UserID, visibilities)

// 	if len(visibilities) == 0 {
// 		return []*models.CalendarEvents{}, nil
// 	}

// 	var events []*models.CalendarEvents
// 	// 💡 Preload("Todos")를 제거하여 Todo 조인을 막습니다.
// 	if err := r.DB.WithContext(ctx).
// 		Where("user_id = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", UserID, visibilities, endAt, startAt).
// 		Order("start_at ASC").
// 		Find(&events).Error; err != nil {
// 		return nil, fmt.Errorf("failed to query events without todos by visibility: %w", err)
// 	}

// 	logger.Infof("Found %d calendar events (without todos) for user %s with visibility filter", len(events), UserID)
// 	return events, nil
// }

// // ------------------------------------------
// // 조회 함수 2: 일별 뷰 (Event + Todo, 캐시 미지원)
// // ------------------------------------------

// // FindCalendarsWithTodos: 특정 기간 동안의 Event와 연결된 Todo를 함께 조회합니다.
// // CalendarService의 GetMyCalendarDailyData/GetUserCalendarDailyData에서 사용됩니다.
// func (r *CalendarEventsRepository) FindCalendarsWithTodos(
// 	ctx context.Context,
// 	UserID uuid.UUID,
// 	visibilities []string,
// 	startAt, endAt time.Time,
// ) ([]*models.CalendarEvents, error) {
// 	logger.Infof("Fetching calendars (with todos) for user=%s with visibilities=%v", UserID, visibilities)

// 	if len(visibilities) == 0 {
// 		return []*models.CalendarEvents{}, nil
// 	}

// 	var events []*models.CalendarEvents
// 	// 💡 Preload("Todos")를 포함하여 Todo를 함께 조회합니다.
// 	if err := r.DB.WithContext(ctx).
// 		Where("user_id = ? AND visibility IN ? AND start_at < ? AND end_at >= ?", UserID, visibilities, endAt, startAt).
// 		Order("start_at ASC").
// 		Preload("Todos").
// 		Find(&events).Error; err != nil {
// 		return nil, fmt.Errorf("failed to query calendars with todos by visibility: %w", err)
// 	}

// 	logger.Infof("Found %d calendar events (with todos) for user %s with visibility filter", len(events), UserID)
// 	return events, nil
// }

// func (r *CalendarEventsRepository) FindEventWithTodosByID(
// 	ctx context.Context,
// 	eventID uuid.UUID,
// ) (*models.CalendarEvents, error) {
// 	var event models.CalendarEvents
// 	// eventID를 사용하여 단일 이벤트를 조회합니다.
// 	if err := r.DB.WithContext(ctx).
// 		Preload("Todos"). // Todo를 함께 로드
// 		First(&event, eventID).Error; err != nil {
// 		// gorm.ErrRecordNotFound 처리를 포함
// 		return nil, fmt.Errorf("failed to query event by ID: %w", err)
// 	}
// 	return &event, nil
// }

// // -------------------------
// // 범위 조회 (visibility 없이, 전체) - 기존 함수 수정 및 유지 (Todos 포함)
// // -------------------------

// func (r *CalendarEventsRepository) FindCalendarsByUser(
// 	ctx context.Context,
// 	UserID uuid.UUID,
// 	startAt, endAt time.Time,
// ) ([]*models.CalendarEvents, error) {
// 	logger.Infof("Fetching ALL calendars for user: %s", UserID)

// 	var events []*models.CalendarEvents
// 	// Todos를 포함하여 조회합니다.
// 	if err := r.DB.WithContext(ctx).
// 		Where("user_id = ? AND start_at < ? AND end_at >= ?", UserID, endAt, startAt).
// 		Order("start_at ASC").
// 		Preload("Todos").
// 		Find(&events).Error; err != nil {
// 		return nil, fmt.Errorf("failed to query calendars: %w", err)
// 	}

// 	logger.Infof("Found %d ALL calendar events for user %s", len(events), UserID)
// 	return events, nil
// }
