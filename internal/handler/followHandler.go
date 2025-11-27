package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

// 라우팅 등록: /me/follows를 RESTful하게 재구성
func (h *FollowHandler) RegisterRoutes(r *gin.Engine) {
	// /me 그룹: 토큰 검증 필수
	me := r.Group("/me")
	me.Use(middleware.AccessTokenAuthMiddleware())

	// /me/follows 그룹: 팔로우/언팔로우/상태 확인 액션
	// 리소스: /me/follows/:target_nickname (인증된 사용자의 팔로우 관계 리소스)
	follows := me.Group("/follows")
	{
		// 1. 팔로우: POST /me/follows/:target_nickname (대상 닉네임을 URI로)
		follows.POST("/:target_nickname", h.Follow)

		// 2. 언팔로우: DELETE /me/follows/:target_nickname
		follows.DELETE("/:target_nickname", h.Unfollow)

		// 3. 팔로우 상태 확인: GET /me/follows/:target_nickname
		// (별도의 /status 대신 GET 요청 자체가 상태 확인 역할)
		follows.GET("/:target_nickname", h.IsFollow)
	}
}

// ---------------------- Handler ----------------------

// follow: POST /me/follows/:target_nickname
func (h *FollowHandler) Follow(c *gin.Context) {
	h.handleFollowAction(c, true)
}

// unfollow: DELETE /me/follows/:target_nickname
func (h *FollowHandler) Unfollow(c *gin.Context) {
	h.handleFollowAction(c, false)
}

// follow/unfollow 처리 공통 함수 (개선: 멱등성 처리)
func (h *FollowHandler) handleFollowAction(c *gin.Context, follow bool) {
	ctx := c.Request.Context()

	// 1. 인증된 사용자 ID (Auth ID) 획득
	authID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 2. 대상 사용자 닉네임/ID 획득
	targetNickname := c.Param("target_nickname")
	if targetNickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_nickname is required"})
		return
	}

	targetID, err := h.ProfileService.GetUserIDByNickname(ctx, targetNickname)
	if err != nil {
		logger.Warnf("failed to get UUID for nickname %s: %v", targetNickname, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 3. 자기 자신 팔로우 방지
	if authID == targetID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot follow or unfollow yourself"})
		return
	}

	// 4. 현재 팔로우 상태 확인
	isFollowing, err := h.FollowService.IsFollow(ctx, authID, targetID)
	if err != nil {
		logger.Errorf("failed to check follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check follow status"})
		return
	}

	// 5. 멱등성 처리: 이미 요청된 상태와 같다면 불필요한 DB 접근을 줄이고 200 OK 반환
	if follow == isFollowing {
		// 이미 팔로우 중인데 팔로우 요청 (POST)
		// 또는 이미 언팔로우 상태인데 언팔로우 요청 (DELETE)
		// -> 요청은 성공적으로 처리된 것으로 간주하고 200 OK 반환 (RESTful 멱등성)

		// 언팔로우 요청의 경우 204 No Content도 흔히 사용됨
		if !follow && !isFollowing {
			c.Status(http.StatusNoContent)
			return
		}

		// 팔로우 요청의 경우 상태 코드 200 유지
		h.respondWithFollowCounts(c, targetNickname, targetID)
		return
	}

	// 6. 팔로우/언팔로우 서비스 호출
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

	// 7. 성공 응답
	h.respondWithFollowCounts(c, targetNickname, targetID)
}

// follow 상태 확인: GET /me/follows/:target_nickname
func (h *FollowHandler) IsFollow(c *gin.Context) {
	ctx := c.Request.Context()

	// 1. 인증된 사용자 ID (Auth ID) 획득
	authID, err := utils.GetUserID(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	targetNickname := c.Param("target_nickname")
	if targetNickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "target_nickname is required"})
		return
	}

	// 2. 대상 사용자 ID (Target ID) 획득
	targetID, err := h.ProfileService.GetUserIDByNickname(ctx, targetNickname)
	if err != nil {
		logger.Warnf("failed to get UUID for nickname %s: %v", targetNickname, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// 3. 팔로우 상태 확인
	isFollowing, err := h.FollowService.IsFollow(ctx, authID, targetID)
	if err != nil {
		logger.Errorf("failed to check follow status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check follow status"})
		return
	}

	// 4. 응답: 팔로우 상태와 대상 닉네임을 반환
	c.JSON(http.StatusOK, gin.H{
		"nickname":     targetNickname,
		"is_following": isFollowing,
	})
}

// 팔로우 수 응답 공통 함수
func (h *FollowHandler) respondWithFollowCounts(c *gin.Context, targetNickname string, targetID uuid.UUID) {
	ctx := c.Request.Context()

	followerCount, followingCount, err := h.ProfileService.GetFollowCounts(ctx, targetID)
	if err != nil {
		logger.Errorf("failed to get follow counts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get follow counts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"profile": gin.H{
			"nickname":        targetNickname,
			"follower_count":  followerCount,
			"following_count": followingCount,
		},
		"message": "follow action successful", // 메시지를 일반화
	})
}
