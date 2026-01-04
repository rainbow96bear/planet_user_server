package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/internal/grpc/client"
	"github.com/rainbow96bear/planet_user_server/internal/models"
	"github.com/rainbow96bear/planet_user_server/internal/repository"
	"github.com/rainbow96bear/planet_utils/pb"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
	"gorm.io/gorm"
)

type TodoServiceInterface interface {
	UpdateTodoStatus(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, isDone bool) (*models.Todo, error)
	FindByID(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*models.Todo, error)
}

type TodoService struct {
	db        *gorm.DB
	TodosRepo *repository.TodosRepository
	Analytics *client.AnalyticsClient
}

func NewTodoService(db *gorm.DB, todosRepo *repository.TodosRepository, analytics *client.AnalyticsClient) *TodoService {
	return &TodoService{
		db:        db,
		TodosRepo: todosRepo,
		Analytics: analytics,
	}
}

// UpdateTodoStatus
func (s *TodoService) UpdateTodoStatus(ctx context.Context, userID uuid.UUID, todoID uuid.UUID, isDone bool) (*models.Todo, error) {
	todo, err := s.TodosRepo.UpdateTodoStatus(ctx, userID, todoID, isDone)
	if err != nil {
		logger.Warnf("todo update failed user=%s todo=%s done=%t err=%v", userID, todoID, isDone, err)
		return nil, err
	}

	logger.Infof("todo status changed user=%s todo=%s done=%t", userID, todoID, isDone)

	// AnalyticsEvent 직접 호출
	if s.Analytics != nil {
		event := &pb.PublishEventRequest{
			EventName:  "todo_completed",
			UserId:     userID.String(),
			OccurredAt: time.Now().Unix(),
			Properties: map[string]string{"todo_id": todoID.String()},
		}
		if !isDone {
			event.EventName = "todo_uncompleted"
		}
		s.Analytics.PublishEvent(ctx, event)
	}

	return todo, nil
}

// FindByID
func (s *TodoService) FindByID(ctx context.Context, userID uuid.UUID, todoID uuid.UUID) (*models.Todo, error) {
	todo, err := s.TodosRepo.FindByID(ctx, todoID)
	if err != nil {
		logger.Warnf("todo find failed user=%s todo=%s err=%v", userID, todoID, err)
		return nil, err
	}
	if todo == nil {
		logger.Infof("todo not found user=%s todo=%s", userID, todoID)
		return nil, fmt.Errorf("todo not found")
	}
	return todo, nil
}
