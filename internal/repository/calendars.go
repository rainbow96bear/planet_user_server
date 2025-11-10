package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/rainbow96bear/planet_utils/model"
)

type CalendarRepository struct {
	DB *sql.DB
}

// -------------------- 조회 --------------------

// GetCalendarByID 단건 조회
func (r *CalendarRepository) GetCalendarByID(ctx context.Context, eventID uint64) (*model.Calendar, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at
		FROM calendar_events
		WHERE id = ?
	`

	cal := &model.Calendar{}
	err := r.DB.QueryRowContext(ctx, query, eventID).Scan(
		&cal.ID,
		&cal.UserUUID,
		&cal.Title,
		&cal.Description,
		&cal.Emoji,
		&cal.StartAt,
		&cal.EndAt,
		&cal.Visibility,
		&cal.ImageURL,
		&cal.CreatedAt,
		&cal.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// Todos 조회
	todos, err := r.getTodosByEventID(ctx, int64(cal.ID))
	if err != nil {
		return nil, err
	}
	cal.Todos = todos
	return cal, nil
}

// GetCalendarsByUserUuid 특정 사용자 전체 일정
func (r *CalendarRepository) GetCalendarsByUserUuid(ctx context.Context, userUUID string) ([]*model.Calendar, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at
		FROM calendar_events
		WHERE user_id = ?
		ORDER BY start_at DESC
	`
	return r.getCalendarList(ctx, query, userUUID)
}

// GetPublicCalendarsByUserUuid 공개 일정 조회
func (r *CalendarRepository) GetPublicCalendarsByUserUuid(ctx context.Context, userUUID string) ([]*model.Calendar, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at
		FROM calendar_events
		WHERE user_id = ? AND visibility = 'public'
		ORDER BY start_at DESC
	`
	return r.getCalendarList(ctx, query, userUUID)
}

// GetAllPublicCalendars 전체 공개 일정
func (r *CalendarRepository) GetAllPublicCalendars(ctx context.Context) ([]*model.Calendar, error) {
	query := `
		SELECT id, user_id, title, description, emoji, start_at, end_at, visibility, image_url, created_at, updated_at
		FROM calendar_events
		WHERE visibility = 'public'
		ORDER BY start_at DESC
	`
	return r.getCalendarList(ctx, query)
}

// -------------------- 생성/수정/삭제 --------------------

// CreateCalendar 새 일정 생성
func (r *CalendarRepository) CreateCalendar(ctx context.Context, cal *model.Calendar) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO calendar_events (user_id, title, description, emoji, start_at, end_at, visibility, image_url)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	result, err := tx.ExecContext(ctx, query,
		cal.UserUUID,
		cal.Title,
		cal.Description,
		cal.Emoji,
		cal.StartAt,
		cal.EndAt,
		cal.Visibility,
		cal.ImageURL,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}
	cal.ID = uint64(id)

	// Todos 삽입
	for _, t := range cal.Todos {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO calendar_todos (event_id, content, done) VALUES (?, ?, ?)",
			cal.ID, t.Content, t.Done,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// UpdateCalendar 일정 수정
func (r *CalendarRepository) UpdateCalendar(ctx context.Context, cal *model.Calendar) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	query := `
		UPDATE calendar_events
		SET title=?, description=?, emoji=?, start_at=?, end_at=?, visibility=?, image_url=?, updated_at=?
		WHERE id=? AND user_id=?
	`
	if _, err := tx.ExecContext(ctx, query,
		cal.Title,
		cal.Description,
		cal.Emoji,
		cal.StartAt,
		cal.EndAt,
		cal.Visibility,
		cal.ImageURL,
		time.Now(),
		cal.ID,
		cal.UserUUID,
	); err != nil {
		tx.Rollback()
		return err
	}

	// Todos 삭제 후 재삽입
	if _, err := tx.ExecContext(ctx, "DELETE FROM calendar_todos WHERE event_id=?", cal.ID); err != nil {
		tx.Rollback()
		return err
	}
	for _, t := range cal.Todos {
		if _, err := tx.ExecContext(ctx,
			"INSERT INTO calendar_todos (event_id, content, done) VALUES (?, ?, ?)",
			cal.ID, t.Content, t.Done,
		); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// DeleteCalendar 일정 삭제
func (r *CalendarRepository) DeleteCalendar(ctx context.Context, eventID uint64) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM calendar_events WHERE id=?", eventID)
	return err
}

// -------------------- Helper --------------------

// Todos 조회
func (r *CalendarRepository) getTodosByEventID(ctx context.Context, eventID int64) ([]model.Todo, error) {
	rows, err := r.DB.QueryContext(ctx,
		"SELECT content, done FROM calendar_todos WHERE event_id=? ORDER BY id ASC",
		eventID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []model.Todo
	for rows.Next() {
		t := model.Todo{}
		if err := rows.Scan(&t.Content, &t.Done); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}
	return todos, nil
}

// 반복되는 calendar_events 조회
func (r *CalendarRepository) getCalendarList(ctx context.Context, query string, args ...interface{}) ([]*model.Calendar, error) {
	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calendars []*model.Calendar
	for rows.Next() {
		c := &model.Calendar{}
		if err := rows.Scan(
			&c.ID,
			&c.UserUUID,
			&c.Title,
			&c.Description,
			&c.Emoji,
			&c.StartAt,
			&c.EndAt,
			&c.Visibility,
			&c.ImageURL,
			&c.CreatedAt,
			&c.UpdatedAt,
		); err != nil {
			return nil, err
		}

		todos, err := r.getTodosByEventID(ctx, int64(c.ID))
		if err != nil {
			return nil, err
		}
		c.Todos = todos
		calendars = append(calendars, c)
	}
	return calendars, nil
}

// 소유자 확인
func (r *CalendarRepository) IsOwnerOfCalendar(ctx context.Context, eventID uint64, userUUID string) (bool, error) {
	query := `SELECT COUNT(*) FROM calendar_events WHERE id = ? AND user_id = ?`
	var count int
	err := r.DB.QueryRowContext(ctx, query, eventID, userUUID).Scan(&count)
	return count > 0, err
}

// FindByNicknameAndVisibility 닉네임 + 공개/친구/개인 조회
func (r *CalendarRepository) FindByNicknameAndVisibility(ctx context.Context, nickname string, visibilityLevels []string) ([]*model.Calendar, error) {
	if len(visibilityLevels) == 0 {
		return []*model.Calendar{}, nil
	}

	placeholders := ""
	args := make([]interface{}, 0, len(visibilityLevels)+1)
	args = append(args, nickname)
	for i := range visibilityLevels {
		if i > 0 {
			placeholders += ", "
		}
		placeholders += "?"
		args = append(args, visibilityLevels[i])
	}

	query := fmt.Sprintf(`
		SELECT ce.id, ce.user_id, ce.title, ce.description, ce.emoji,
		       ce.start_at, ce.end_at, ce.visibility, ce.image_url, ce.created_at, ce.updated_at
		FROM calendar_events AS ce
		JOIN profiles AS p ON p.uuid = ce.user_id
		WHERE p.nickname = ? AND ce.visibility IN (%s)
		ORDER BY ce.start_at DESC
	`, placeholders)

	return r.getCalendarList(ctx, query, args...)
}
