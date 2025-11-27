package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type CalendarHandler struct {
	CalendarService *service.CalendarService
	FollowService   *service.FollowService
	ProfileService  *service.ProfileService
}

func NewCalendarHandler(calendarService *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{CalendarService: calendarService}
}

func (h *CalendarHandler) RegisterRoutes(r *gin.Engine) {
	// 공개용: 다른 사용자가 보는 달력
	// Profile 기반, 공개 여부 확인
	me := r.Group("/me")
	me.Use(middleware.AccessTokenAuthMiddleware())
	{
		me.GET("/calendar", h.GetMyCalendarEvent)
		me.POST("/calendar/events", h.CreateCalendarEvent)
		me.PUT("/calendar/events/:eventId", h.UpdateCalendarEvent)
		me.DELETE("/calendar/events/:eventId", h.DeleteCalendarEvent)
	}

	users := r.Group("/users/:nickname")
	users.GET("/calendar", h.GetUserCalendarEvent)
}

// ---------------------- Handler ----------------------

// 다른 사람 캘린더 조회
func (h *CalendarHandler) GetUserCalendarEvent(c *gin.Context) {
	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	logger.Infof("GetUserCalendar start, nickname=%s", nickname)
	defer logger.Infof("GetUserCalendar end, nickname=%s", nickname)

	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	data, err := h.CalendarService.GetUserCalendarData(ctx, nickname, authID, year, month)
	if err != nil {
		logger.Errorf("GetUserCalendar failed, nickname=%s, err=%v", nickname, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendars"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// 내 캘린더 조회
func (h *CalendarHandler) GetMyCalendarEvent(c *gin.Context) {
	ctx := c.Request.Context()
	UserID, _ := utils.GetUserID(c)
	logger.Infof("GetMyCalendar start, UserID=%s", UserID)
	defer logger.Infof("GetMyCalendar end, UserID=%s", UserID)

	year, month, err := parseYearMonth(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	data, err := h.CalendarService.GetMyCalendarData(ctx, UserID, year, month)
	if err != nil {
		logger.Errorf("GetMyCalendar failed, UserID=%s, err=%v", UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get calendars"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// 캘린더 생성
func (h *CalendarHandler) CreateCalendarEvent(c *gin.Context) {
	ctx := c.Request.Context()
	UserID, _ := utils.GetUserID(c)
	logger.Infof("CreateCalendar start, UserID=%s", UserID)
	defer logger.Infof("CreateCalendar end, UserID=%s", UserID)

	var req dto.CalendarCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("CreateCalendar invalid request, UserID=%s, err=%v", UserID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "detail": err.Error()})
		return
	}

	calendar := dto.ToCalendarModelFromCreate(&req, UserID)
	if err := h.CalendarService.CreateCalendarEvent(ctx, calendar); err != nil {
		logger.Errorf("CreateCalendar failed, UserID=%s, err=%v", UserID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create calendar"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "calendar created",
		"event":   dto.ToCalendarInfo(calendar),
	})
}

// 캘린더 업데이트
func (h *CalendarHandler) UpdateCalendarEvent(c *gin.Context) {
	ctx := c.Request.Context()
	eventIDStr := c.Param("eventId")
	UserID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}
	logger.Infof("UpdateCalendar start, UserID=%s, eventID=%s", UserID, eventIDStr)
	defer logger.Infof("UpdateCalendar end, UserID=%s, eventID=%s", UserID, eventIDStr)

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	var req dto.CalendarUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "detail": err.Error()})
		return
	}

	if err := h.CalendarService.UpdateCalendarEvent(ctx, UserID, eventID, &req); err != nil {
		logger.Errorf("UpdateCalendar failed, UserID=%s, eventID=%s, err=%v", UserID, eventID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update calendar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "calendar updated successfully"})
}

// 캘린더 삭제
func (h *CalendarHandler) DeleteCalendarEvent(c *gin.Context) {
	ctx := c.Request.Context()
	eventIDStr := c.Param("eventId")
	UserID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}
	logger.Infof("DeleteCalendar start, UserID=%s, eventID=%s", UserID, eventIDStr)
	defer logger.Infof("DeleteCalendar end, UserID=%s, eventID=%s", UserID, eventIDStr)

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	if err := h.CalendarService.DeleteCalendarEvent(ctx, UserID, eventID); err != nil {
		logger.Errorf("DeleteCalendar failed, UserID=%s, eventID=%s, err=%v", UserID, eventID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete calendar"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "calendar deleted successfully"})
}

// ---------------------- Helper ----------------------
func parseYearMonth(c *gin.Context) (int, int, error) {
	yearStr := c.Query("year")
	monthStr := c.Query("month")
	if yearStr == "" || monthStr == "" {
		return 0, 0, fmt.Errorf("year and month are required")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year")
	}
	month, err := strconv.Atoi(monthStr)
	if err != nil || month < 1 || month > 12 {
		return 0, 0, fmt.Errorf("invalid month")
	}

	return year, month, nil
}
