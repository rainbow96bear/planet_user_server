package resolver

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/graph/model"
	"github.com/rainbow96bear/planet_user_server/internal/mapper"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

//
// =======================
// Utils
// =======================
//

func monthRange(year int32, month int32) (time.Time, time.Time, error) {
	if month < 1 || month > 12 {
		return time.Time{}, time.Time{}, errors.New("invalid month")
	}

	loc := time.UTC
	from := time.Date(int(year), time.Month(month), 1, 0, 0, 0, 0, loc)
	to := from.AddDate(0, 1, 0)

	return from, to, nil
}

//
// =======================
// Mutation
// =======================
//

func (r *mutationResolver) CreateCalendarEvent(
	ctx context.Context,
	input model.CreateCalendarInput,
) (*model.Calendar, error) {

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	calendar := dto.ToCalendarModel(input, userID)

	created, err := r.CalendarService.Create(ctx, calendar)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQL(created), nil
}

func (r *mutationResolver) UpdateCalendarEvent(
	ctx context.Context,
	eventID string,
	input model.UpdateCalendarInput,
) (*model.Calendar, error) {

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, errors.New("invalid event id")
	}

	event, err := r.CalendarService.Update(
		ctx,
		userID,
		eventUUID,
		input,
	)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQL(event), nil
}

func (r *mutationResolver) DeleteCalendarEvent(
	ctx context.Context,
	eventID string,
) (bool, error) {

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return false, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return false, errors.New("unauthorized")
	}

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return false, errors.New("invalid event id")
	}

	if err := r.CalendarService.Delete(ctx, userID, eventUUID); err != nil {
		return false, err
	}

	return true, nil
}

//
// =======================
// Query
// =======================
//

// 로그인 사용자 월별 일정 (todo 없음)
func (r *queryResolver) MyCalendarEvents(
	ctx context.Context,
	year int32,
	month int32,
) ([]*model.Calendar, error) {

	from, to, err := monthRange(year, month)
	if err != nil {
		return nil, err
	}

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	logger.Infof("MyCalendarEvents user=%s from=%s to=%s",
		userID, from.Format(time.RFC3339), to.Format(time.RFC3339))

	events, err := r.CalendarService.GetMyEventsByPeriod(
		ctx,
		userID,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQLList(events), nil
}

// 로그인 사용자 단건 일정 (todo 포함)
func (r *queryResolver) MyCalendarEvent(
	ctx context.Context,
	eventID string,
) (*model.Calendar, error) {

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	eventUUID, err := uuid.Parse(eventID)
	if err != nil {
		return nil, errors.New("invalid event id")
	}

	event, err := r.CalendarService.GetDetailWithTodos(
		ctx,
		userID,
		eventUUID,
	)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQL(event), nil
}

// 로그인 사용자 일별 일정 (todo 포함)
func (r *queryResolver) MyCalendarEventsByDate(
	ctx context.Context,
	date time.Time,
) ([]*model.Calendar, error) {

	token, err := middleware.ExtractAccessToken(ctx)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	userID, err := utils.GetUserID(token)
	if err != nil {
		return nil, errors.New("unauthorized")
	}

	events, err := r.CalendarService.GetMyEventsByDateWithTodos(
		ctx,
		userID,
		date,
	)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQLList(events), nil
}

// 다른 사용자 월별 일정 (비로그인 가능)
func (r *queryResolver) UserCalendarEvents(
	ctx context.Context,
	userID string,
	year int32,
	month int32,
) ([]*model.Calendar, error) {

	from, to, err := monthRange(year, month)
	if err != nil {
		return nil, err
	}

	targetUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// viewer는 optional
	var viewerUUID *uuid.UUID
	if token, err := middleware.ExtractAccessToken(ctx); err == nil {
		if uid, err := utils.GetUserID(token); err == nil {
			viewerUUID = &uid
		}
	}

	logger.Infof(
		"UserCalendarEvents target=%s viewer=%v from=%s to=%s",
		targetUUID,
		viewerUUID,
		from.Format(time.RFC3339),
		to.Format(time.RFC3339),
	)

	events, err := r.CalendarService.GetUserEventsByPeriod(
		ctx,
		viewerUUID,
		targetUUID,
		from,
		to,
	)
	if err != nil {
		return nil, err
	}

	return mapper.ToCalendarGraphQLList(events), nil
}
