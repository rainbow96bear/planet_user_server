package utils

import (
	"fmt"
	"net/http"
	"reflect"

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

func StructToUpdateMap(input any) map[string]any {
	result := make(map[string]any)

	v := reflect.ValueOf(input).Elem() // *struct → struct
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// json 태그 읽기
		tag := t.Field(i).Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}
		name := tag
		// json:"nickname,omitempty" → nickname 추출
		if commaIdx := len(name) - len(",omitempty"); commaIdx > 0 && name[commaIdx:] == ",omitempty" {
			name = name[:commaIdx]
		}

		// nil 포인터는 스킵
		if field.Kind() == reflect.Pointer && field.IsNil() {
			continue
		}

		// 포인터면 값 가져오기, 아니면 그냥 저장
		if field.Kind() == reflect.Pointer {
			result[name] = field.Elem().Interface()
		} else {
			result[name] = field.Interface()
		}
	}

	return result
}
