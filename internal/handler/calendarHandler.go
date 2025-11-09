package handler

import (
	"encoding/json"
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
	FollowService   *service.FollowService
	// 추후 이미지 업로드 서비스 추가
	// ImageUploader dto.ImageUploader
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

	// FormData 바인딩
	var req dto.CalendarCreateRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Errorf("form bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// todos 파싱
	var todos []dto.TodoItem
	if req.Todos != "" {
		if err := json.Unmarshal([]byte(req.Todos), &todos); err != nil {
			logger.Errorf("todos parse error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todos format"})
			return
		}
	}

	// 이미지 업로드 처리 (추후 활성화)
	// var imageURL string
	// if req.HasImage() {
	// 	imageURL, err = h.ImageUploader.Upload(req.Image)
	// 	if err != nil {
	// 		logger.Errorf("image upload error: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
	// 		return
	// 	}
	// }

	calendar := dto.ToCalendarModel(&req, userUUID, todos)

	if err := h.CalendarService.CreateCalendar(ctx, calendar); err != nil {
		logger.Errorf("create calendar error: %v", err)
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

	// FormData 바인딩 (JSON 대신 Form으로 변경)
	var req dto.CalendarUpdateRequest
	if err := c.ShouldBind(&req); err != nil {
		logger.Errorf("form bind error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form data"})
		return
	}

	// todos 파싱
	var todos []dto.TodoItem
	if req.Todos != "" {
		if err := json.Unmarshal([]byte(req.Todos), &todos); err != nil {
			logger.Errorf("todos parse error: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid todos format"})
			return
		}
	}

	// 기존 일정 조회 (이미지 URL 유지를 위해)
	existingCalendar, err := h.CalendarService.GetCalendarByID(ctx, userUUID, eventId)
	if err != nil {
		logger.Errorf("get existing calendar error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "calendar not found"})
		return
	}

	imageURL := existingCalendar.ImageURL

	// 이미지 업로드/삭제 처리 (추후 활성화)
	// if req.DeleteImage {
	// 	if existingCalendar.ImageURL != "" {
	// 		_ = h.ImageUploader.Delete(existingCalendar.ImageURL)
	// 	}
	// 	imageURL = ""
	// } else if req.Image != nil && req.Image.Size > 0 {
	// 	// 기존 이미지 삭제
	// 	if existingCalendar.ImageURL != "" {
	// 		_ = h.ImageUploader.Delete(existingCalendar.ImageURL)
	// 	}
	// 	// 새 이미지 업로드
	// 	imageURL, err = h.ImageUploader.Upload(req.Image)
	// 	if err != nil {
	// 		logger.Errorf("image upload error: %v", err)
	// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload image"})
	// 		return
	// 	}
	// }

	calendar := dto.ToCalendarUpdateModel(&req, userUUID, eventId, todos, imageURL)

	if err := h.CalendarService.UpdateCalendar(ctx, calendar); err != nil {
		logger.Errorf("update calendar error: %v", err)
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

	// 삭제 전 일정 조회 (이미지 삭제를 위해, 추후 사용)
	// calendar, err := h.CalendarService.GetCalendarByID(ctx, userUUID, eventId)
	// if err != nil {
	// 	logger.Errorf("get calendar error: %v", err)
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "calendar not found"})
	// 	return
	// }

	// 이미지 삭제 (추후 활성화)
	// if calendar.ImageURL != "" {
	// 	_ = h.ImageUploader.Delete(calendar.ImageURL)
	// }

	if err := h.CalendarService.DeleteCalendar(ctx, userUUID, eventId); err != nil {
		logger.Errorf("delete calendar error: %v", err)
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
	isFollowing, err := h.FollowService.IsFollow(ctx, userUUID, nickname)
	if err != nil {
		logger.Warnf("follow status check failed: %v", err)
	}

	var calendars []*dto.CalendarInfo
	if isFollowing {
		calendars, err = h.CalendarService.GetPublicCalendarByNickname(ctx, nickname)
		if err != nil {
			logger.Errorf("get user calendar error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		calendars, err = h.CalendarService.GetPublicCalendarByNickname(ctx, nickname)
		if err != nil {
			logger.Errorf("get user calendar error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
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
		logger.Errorf("get all public calendars error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load public calendars"})
		return
	}

	c.JSON(http.StatusOK, calendars)
}
