package utils

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/rainbow96bear/planet_user_server/config"
)

func GetUuidByAccessToken(c *gin.Context) (string, error) {
	tokenStr, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing access token"})
		return "", err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JWT_SECRET_KEY), nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return "", fmt.Errorf("fail to claims")
	}

	userUuid, ok := claims["userUuid"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user uuid in token"})
		return "", fmt.Errorf("fail to get uuid from jwt")
	}

	return userUuid, nil

}
