package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/models"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type TodosRepository struct {
	DB *gorm.DB
}

func (r *TodosRepository) BeginTx(ctx context.Context) (*gorm.DB, error) {
	logger.Infof("starting transaction for TodosRepository")
	tx := r.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		logger.Errorf("failed to start transaction: %v", tx.Error)
		return nil, tx.Error
	}
	logger.Infof("transaction started successfully")
	return tx, nil
}

// -------------------------
// 단일 Todo 생성
// -------------------------
func (r *TodosRepository) CreateTodo(ctx context.Context, todo *models.Todos) error {
	logger.Infof("Creating todo for event: %s", todo.EventID)

	if err := r.DB.WithContext(ctx).Create(todo).Error; err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}

	logger.Infof("Successfully created todo: %s", todo.ID)
	return nil
}

// -------------------------
// Todo 상태 업데이트
// -------------------------
func (r *TodosRepository) UpdateTodoStatus(ctx context.Context, todoID uuid.UUID, isDone bool) error {
	logger.Infof("Updating todo status: %s to %v", todoID, isDone)

	if err := r.DB.WithContext(ctx).
		Model(&models.Todos{}).
		Where("id = ?", todoID).
		Update("is_done", isDone).Error; err != nil {
		return fmt.Errorf("failed to update todo status: %w", err)
	}

	logger.Infof("Todo %s status updated to %v", todoID, isDone)
	return nil
}

// -------------------------
// EventID 기반 Todo 조회
// -------------------------
func (r *TodosRepository) FindTodosByEventID(ctx context.Context, eventID uuid.UUID) ([]*models.Todos, error) {
	var todos []*models.Todos
	if err := r.DB.WithContext(ctx).Where("event_id = ?", eventID).Find(&todos).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch todos: %w", err)
	}
	logger.Infof("Found %d todos for event %s", len(todos), eventID)
	return todos, nil
}

// -------------------------
// 단일 Todo 삭제
// -------------------------
func (r *TodosRepository) DeleteTodo(ctx context.Context, todoID uuid.UUID) error {
	if err := r.DB.WithContext(ctx).Delete(&models.Todos{}, "id = ?", todoID).Error; err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}
	logger.Infof("Todo %s deleted", todoID)
	return nil
}
