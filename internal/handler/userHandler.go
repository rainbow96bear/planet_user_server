package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type UserHandler struct {
	ProfileService *service.ProfileService
	FollowService  *service.FollowService
}

func NewUserHandler(profileService *service.ProfileService, followService *service.FollowService) *UserHandler {
	return &UserHandler{
		ProfileService: profileService,
		FollowService:  followService,
	}
}

func (h *UserHandler) RegisterRoutes(r *gin.Engine) {
	usersGroup := r.Group("/users")
	usersGroup.Use(middleware.AuthMiddleware())
	{
		usersGroup.POST("/:nickname/follow", h.Follow)
		usersGroup.DELETE("/:nickname/follow", h.Unfollow)
		usersGroup.GET("/:nickname/follow-status", h.IsFollow)
	}
}

func (h *UserHandler) Follow(c *gin.Context) {
	logger.Infof("start to follow")
	defer logger.Infof("end to follow")

	ctx := c.Request.Context()
	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followeeNickname := c.Param("nickname")

	// followeeNickname을 uuid로 변환
	followeeUuid, err := h.ProfileService.GetUserUuidByNickname(ctx, followeeNickname)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 이미 follow 중인지 체크
	isFollow, err := h.FollowService.IsFollow(ctx, authUuid, followeeUuid)
	if err != nil {
		logger.Errorf("fail to check isfollow ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isFollow {
		logger.Warnf("Already follow: %s -> %s", authUuid, followeeUuid)
		c.JSON(http.StatusConflict, gin.H{"error": "already following this user"})
		return
	}

	// follow table에 추가
	err = h.FollowService.Follow(ctx, authUuid, followeeUuid)
	if err != nil {
		logger.Errorf("fail to check isfollow ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	followerCounts, followeeCounts, err := h.ProfileService.GetFollowCounts(ctx, followeeUuid)
	if err != nil {
		logger.Errorf("fail to get follow counts ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 응답
	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"nickname":        followeeNickname,
			"follow_count":    followerCounts,
			"following_count": followeeCounts,
		},
		"message": "successfully followed",
	})
}

func (h *UserHandler) Unfollow(c *gin.Context) {
	logger.Infof("start to unfollow")
	defer logger.Infof("end to unfollow")

	ctx := c.Request.Context()
	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followeeNickname := c.Param("nickname")

	// followeeNickname을 uuid로 변환
	followeeUuid, err := h.ProfileService.GetUserUuidByNickname(ctx, followeeNickname)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 이미 follow 중인지 체크
	isFollow, err := h.FollowService.IsFollow(ctx, authUuid, followeeUuid)
	if err != nil {
		logger.Errorf("fail to check isfollow ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isFollow {
		logger.Warnf("Already unfollow: %s -> %s", authUuid, followeeUuid)
		c.JSON(http.StatusConflict, gin.H{"error": "already unfollowing this user"})
		return
	}

	// follow table에 추가
	err = h.FollowService.Unfollow(ctx, authUuid, followeeUuid)
	if err != nil {
		logger.Errorf("fail to check isfollow ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	followerCounts, followeeCounts, err := h.ProfileService.GetFollowCounts(ctx, followeeUuid)
	if err != nil {
		logger.Errorf("fail to get follow counts ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 응답
	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"nickname":        followeeNickname,
			"follow_count":    followerCounts,
			"following_count": followeeCounts,
		},
		"message": "successfully followed",
	})
}

func (h *UserHandler) IsFollow(c *gin.Context) {
	logger.Infof("start to check follow status")
	defer logger.Infof("end to check follow status")

	ctx := c.Request.Context()
	authUuid, err := utils.GetUserUuid(c)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	followeeNickname := c.Param("nickname")

	// followeeNickname을 uuid로 변환
	followeeUuid, err := h.ProfileService.GetUserUuidByNickname(ctx, followeeNickname)
	if err != nil {
		logger.Errorf(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 이미 follow 중인지 체크
	isFollow, err := h.FollowService.IsFollow(ctx, authUuid, followeeUuid)
	if err != nil {
		logger.Errorf("fail to check isfollow ERR[%s]", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nickname":     followeeNickname,
		"is_following": isFollow,
	})
}
