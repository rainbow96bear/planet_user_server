package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rainbow96bear/planet_user_server/internal/service"
	"github.com/rainbow96bear/planet_user_server/middleware"
	"github.com/rainbow96bear/planet_user_server/utils"
	"github.com/rainbow96bear/planet_utils/pkg/logger"
)

type FollowHandler struct {
	ProfileService *service.ProfileService
	FollowService  *service.FollowService
}

func NewFollowHandler(profileService *service.ProfileService, followService *service.FollowService) *FollowHandler {
	return &FollowHandler{
		ProfileService: profileService,
		FollowService:  followService,
	}
}

func (h *FollowHandler) RegisterRoutes(r *gin.Engine) {
	usersGroup := r.Group("/follow")
	usersGroup.Use(middleware.AuthMiddleware())
	{
		usersGroup.POST("/:nickname", h.Follow)
		usersGroup.DELETE("/:nickname", h.Unfollow)
		usersGroup.GET("/:nickname/status", h.IsFollow)
	}
}

// ---------------------- Handler ----------------------

// follow
func (h *FollowHandler) Follow(c *gin.Context) {
	h.handleFollowAction(c, true)
}

// unfollow
func (h *FollowHandler) Unfollow(c *gin.Context) {
	h.handleFollowAction(c, false)
}

// follow/unfollow 처리 공통 함수
func (h *FollowHandler) handleFollowAction(c *gin.Context, follow bool) {
	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)
	if err != nil {
		logger.Errorf("failed to get auth ID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetNickname := c.Param("nickname")
	if targetNickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	targetID, err := h.ProfileService.GetUserIDByNickname(ctx, targetNickname)
	if err != nil {
		logger.Errorf("failed to get UUID for nickname %s: %v", targetNickname, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve user"})
		return
	}

	isFollowing, err := h.FollowService.IsFollow(ctx, authID, targetID)
	if err != nil {
		logger.Errorf("failed to check follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check follow status"})
		return
	}

	// 상태 확인 후 처리
	if follow && isFollowing {
		c.JSON(http.StatusConflict, gin.H{"error": "already following this user"})
		return
	}
	if !follow && !isFollowing {
		c.JSON(http.StatusConflict, gin.H{"error": "already unfollowing this user"})
		return
	}

	if follow {
		if err := h.FollowService.Follow(ctx, authID, targetID); err != nil {
			logger.Errorf("failed to follow user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
			return
		}
	} else {
		if err := h.FollowService.Unfollow(ctx, authID, targetID); err != nil {
			logger.Errorf("failed to unfollow user: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
			return
		}
	}

	followerCount, followingCount, err := h.ProfileService.GetFollowCounts(ctx, targetID)
	if err != nil {
		logger.Errorf("failed to get follow counts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follow counts"})
		return
	}

	action := "followed"
	if !follow {
		action = "unfollowed"
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"nickname":        targetNickname,
			"follower_count":  followerCount,
			"following_count": followingCount,
		},
		"message": "successfully " + action,
	})
}

// follow 상태 확인
func (h *FollowHandler) IsFollow(c *gin.Context) {
	ctx := c.Request.Context()
	authID, err := utils.GetUserID(c)
	if err != nil {
		logger.Errorf("failed to get auth UUID: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetNickname := c.Param("nickname")
	if targetNickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	targetID, err := h.ProfileService.GetUserIDByNickname(ctx, targetNickname)
	if err != nil {
		logger.Errorf("failed to get UUID for nickname %s: %v", targetNickname, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to resolve user"})
		return
	}

	isFollowing, err := h.FollowService.IsFollow(ctx, authID, targetID)
	if err != nil {
		logger.Errorf("failed to check follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check follow status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"nickname":     targetNickname,
		"is_following": isFollowing,
	})
}
