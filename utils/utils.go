package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserID(c *gin.Context) (uuid.UUID, error) {
	authUuidValue, exists := c.Get("user_uuid")
	if !exists {
		err := fmt.Errorf("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return uuid.Nil, err
	}

	authIDStr, ok := authUuidValue.(string)
	if !ok {
		err := fmt.Errorf("invalid user uuid type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return uuid.Nil, err
	}
	authID, err := uuid.Parse(authIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid eventId"})
		return uuid.Nil, err
	}
	return authID, nil
}

func GetUserNickname(c *gin.Context) (string, error) {
	authNicknameValue, exists := c.Get("nickname")
	if !exists {
		err := fmt.Errorf("unauthorized")
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return "", err
	}

	authNickname, ok := authNicknameValue.(string)
	if !ok {
		err := fmt.Errorf("invalid user nickname type")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return "", err
	}
	return authNickname, nil
}
