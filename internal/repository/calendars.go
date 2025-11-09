package repository

import (
	"context"
	"database/sql"

	"github.com/rainbow96bear/planet_user_server/dto"
)

type CalendarRepository struct {
	DB *sql.DB
}

// =====================
// Calendar 조회
// =====================

// GetCalendarByID 특정 일정 조회 (단건)
func (r *CalendarRepository) GetCalendarByID(ctx context.Context, eventId int64) (*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE id = ?
	`
	c := &dto.CalendarInfo{}

	err := r.DB.QueryRowContext(ctx, query, eventId).Scan(
		&c.EventID,
		&c.UserUUID,
		&c.Title,
		&c.Description,
		&c.Emoji,
		&c.StartAt,
		&c.EndAt,
		&c.Visibility,
		&c.ImageURL,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// 할 일 조회
	todos, err := r.getTodosByEventID(ctx, eventId)
	if err != nil {
		return nil, err
	}
	c.Todos = todos

	return c, nil
}

// GetCalendarsByUserUuid 내 일정 목록
func (r *CalendarRepository) GetCalendarsByUserUuid(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE user_id = ?
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(
			&c.EventID,
			&c.UserUUID,
			&c.Title,
			&c.Description,
			&c.Emoji,
			&c.StartAt,
			&c.EndAt,
			&c.Visibility,
			&c.ImageURL,
		); err != nil {
			return nil, err
		}

		todos, err := r.getTodosByEventID(ctx, c.EventID)
		if err != nil {
			return nil, err
		}
		c.Todos = todos

		calendars = append(calendars, c)
	}
	return calendars, nil
}

// =====================
// Calendar 생성 / 수정 / 삭제
// =====================

// CreateCalendar 일정 생성
func (r *CalendarRepository) CreateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	// calendar_events 삽입
	query := `
		INSERT INTO calendar_events (user_id, title, description, emoji, start_at, end_at, visibility, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, query,
		calendar.UserUUID,
		calendar.Title,
		calendar.Description,
		calendar.Emoji,
		calendar.StartAt,
		calendar.EndAt,
		calendar.Visibility,
		calendar.ImageURL,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	eventID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	calendar.EventID = eventID

	// todos 삽입
	for _, t := range calendar.Todos {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO calendar_todos (event_id, content, done) VALUES (?, ?, ?)",
			eventID, t.Text, t.Completed,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// UpdateCalendar 일정 수정
func (r *CalendarRepository) UpdateCalendar(ctx context.Context, calendar *dto.CalendarInfo) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		UPDATE calendar_events
		SET title=?, description=?, emoji=?, start_at=?, end_at=?, visibility=?, image_url=?
		WHERE id=? AND user_id=?
	`
	if _, err := tx.ExecContext(ctx, query,
		calendar.Title,
		calendar.Description,
		calendar.Emoji,
		calendar.StartAt,
		calendar.EndAt,
		calendar.Visibility,
		calendar.ImageURL,
		calendar.EventID,
		calendar.UserUUID,
	); err != nil {
		tx.Rollback()
		return err
	}

	// todos 삭제 후 재삽입 (간단한 방법)
	if _, err := tx.ExecContext(ctx, "DELETE FROM calendar_todos WHERE event_id=?", calendar.EventID); err != nil {
		tx.Rollback()
		return err
	}
	for _, t := range calendar.Todos {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO calendar_todos (event_id, content, done) VALUES (?, ?, ?)",
			calendar.EventID, t.Text, t.Completed,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeleteCalendar 일정 삭제
func (r *CalendarRepository) DeleteCalendar(ctx context.Context, eventId int64) error {
	// calendar_todos는 foreign key cascade 처리됨
	_, err := r.DB.ExecContext(ctx, "DELETE FROM calendar_events WHERE id=?", eventId)
	return err
}

// =====================
// Helper: Todos 조회
// =====================
func (r *CalendarRepository) getTodosByEventID(ctx context.Context, eventID int64) ([]dto.TodoItem, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT content, done FROM calendar_todos WHERE event_id=? ORDER BY id ASC",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []dto.TodoItem
	for rows.Next() {
		t := dto.TodoItem{}
		if err := rows.Scan(&t.Text, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// =====================
// 소유자 확인 / 공개 일정
// =====================
func (r *CalendarRepository) IsOwnerOfCalendar(ctx context.Context, eventId int64, userUuid string) (bool, error) {
	query := `SELECT COUNT(*) FROM calendar_events WHERE id = ? AND user_id = ?`
	var count int
	err := r.DB.QueryRowContext(ctx, query, eventId, userUuid).Scan(&count)
	return count > 0, err
}

func (r *CalendarRepository) GetPublicCalendarsByUserUuid(ctx context.Context, userUuid string) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE user_id = ? AND visibility = 'public'
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query, userUuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(
			&c.EventID,
			&c.UserUUID,
			&c.Title,
			&c.Description,
			&c.Emoji,
			&c.StartAt,
			&c.EndAt,
			&c.Visibility,
			&c.ImageURL,
		); err != nil {
			return nil, err
		}

		todos, err := r.getTodosByEventID(ctx, c.EventID)
		if err != nil {
			return nil, err
		}
		c.Todos = todos

		calendars = append(calendars, c)
	}
	return calendars, nil
}

func (r *CalendarRepository) GetAllPublicCalendars(ctx context.Context) ([]*dto.CalendarInfo, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url
		FROM calendar_events
		WHERE visibility = 'public'
		ORDER BY start_at DESC
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*dto.CalendarInfo
	for rows.Next() {
		c := &dto.CalendarInfo{}
		if err := rows.Scan(
			&c.EventID,
			&c.UserUUID,
			&c.Title,
			&c.Description,
			&c.Emoji,
			&c.StartAt,
			&c.EndAt,
			&c.Visibility,
			&c.ImageURL,
		); err != nil {
			return nil, err
		}

		todos, err := r.getTodosByEventID(ctx, c.EventID)
		if err != nil {
			return nil, err
		}
		c.Todos = todos

		calendars = append(calendars, c)
	}
	return calendars, nil
}
