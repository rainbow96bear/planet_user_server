package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/dto"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type CalendarHandler struct {
	CalendarService *service.CalendarService
}

func NewCalendarHandler(calendarService *service.CalendarService) *CalendarHandler {
	return &CalendarHandler{
		CalendarService: calendarService,
	}
}

func (h *CalendarHandler) RegisterRoutes(r *gin.Engine) {
	calendarGroup := r.Group("/calendar")

	// 공개용 (다른 사람 일정 조회)
	calendarGroup.GET("/user/:nickname", h.GetUserCalendar)
	calendarGroup.GET("/user", h.GetAllPublicCalendars)

	// 인증 필요
	calendarGroup.Use(middleware.AuthMiddleware())
	{
		calendarGroup.GET("", h.GetMyCalendar)
		calendarGroup.POST("", h.CreateCalendar)
		calendarGroup.PUT("/:eventId", h.UpdateCalendar)
		calendarGroup.DELETE("/:eventId", h.DeleteCalendar)
	}
}

// GET /calendar
func (h *CalendarHandler) GetMyCalendar(c *gin.Context) {
	logger.Infof("start to get my calendar")
	defer logger.Infof("end to get my calendar")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserUuid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	calendars, err := h.CalendarService.GetUserCalendar(ctx, userUUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load calendar"})
		return
	}

	c.JSON(http.StatusOK, calendars)
}

// POST /calendar
func (h *CalendarHandler) CreateCalendar(c *gin.Context) {
	logger.Infof("start to create calendar")
	defer logger.Infof("end to create calendar")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserUuid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	var req dto.CalendarCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	calendar := dto.ToCalendarModel(&req, userUUID)
	if err := h.CalendarService.CreateCalendar(ctx, calendar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, calendar)
}

// PUT /calendar/:eventId
func (h *CalendarHandler) UpdateCalendar(c *gin.Context) {
	logger.Infof("start to update calendar")
	defer logger.Infof("end to update calendar")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserUuid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	eventIdStr := c.Param("eventId")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	var req dto.CalendarUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	calendar := dto.ToCalendarUpdateModel(&req, userUUID, eventId)
	if err := h.CalendarService.UpdateCalendar(ctx, calendar); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, calendar)
}

// DELETE /calendar/:eventId
func (h *CalendarHandler) DeleteCalendar(c *gin.Context) {
	logger.Infof("start to delete calendar")
	defer logger.Infof("end to delete calendar")

	ctx := c.Request.Context()
	userUUID, err := utils.GetUserUuid(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	eventIdStr := c.Param("eventId")
	eventId, err := strconv.ParseInt(eventIdStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return
	}

	if err := h.CalendarService.DeleteCalendar(ctx, userUUID, eventId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true})
}

// GET /calendar/user/:nickname
func (h *CalendarHandler) GetUserCalendar(c *gin.Context) {
	logger.Infof("start to get user calendar")
	defer logger.Infof("end to get user calendar")

	ctx := c.Request.Context()
	nickname := c.Param("nickname")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	calendars, err := h.CalendarService.GetPublicCalendarByNickname(ctx, nickname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, calendars)
}

// GET /calendar/user
func (h *CalendarHandler) GetAllPublicCalendars(c *gin.Context) {
	logger.Infof("start to get all public calendars")
	defer logger.Infof("end to get all public calendars")

	ctx := c.Request.Context()
	calendars, err := h.CalendarService.GetAllPublicCalendars(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load public calendars"})
		return
	}

	c.JSON(http.StatusOK, calendars)
}
