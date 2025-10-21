package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/rainbow96bear/planet_proto"
	"github.com/rainbow96bear/planet_user_server/config"
	"github.com/rainbow96bear/planet_user_server/grpc_client"
	"github.com/rainbow96bear/planet_user_server/logger"
)

// 프로필 조회
func GetProfile(c *gin.Context) {
	nickname := c.Param("userNickName")
	if nickname == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "nickname is required"})
		return
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}

	res, err := dbClient.ReqGetUserInfoByNickname(&pb.UserInfo{Nickname: nickname})
	if err != nil {
		logger.Warnf("failed to get profile: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"userUuid":     res.UserUuid,
		"nickname":     res.Nickname,
		"profileImage": res.ProfileImage,
		"role":         res.Role,
	})
}

// 프로필 수정 (로그인 사용자 본인만)
func UpdateProfile(c *gin.Context) {
	userUuid, exists := c.Get("userUuid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req struct {
		Nickname     string `json:"nickname"`
		ProfileImage string `json:"profile_image"`
		Bio          string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	dbClient, err := grpc_client.NewDBClient(config.DB_GRPC_SERVER_ADDR)
	if err != nil {
		logger.Errorf("fail to create dbClient ERR[%s]", err.Error())
		return
	}

	res, err := dbClient.ReqUpdateUserProfile(&pb.UserInfo{
		UserUuid:     userUuid.(string),
		Nickname:     req.Nickname,
		ProfileImage: req.ProfileImage,
		Bio:          req.Bio,
	})

	if err != nil {
		logger.Warnf("failed to update profile: %+v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "profile updated",
		"user":    res,
	})
}
