package planet_err

import (
	"errors"

	"gorm.io/gorm"
)

// 리소스 관련 일반 오류
var ErrNotFound = errors.New("resource not found")
var ErrAlreadyExists = errors.New("resource already exists")

// 닉네임 중복 오류
var ErrNicknameDuplicate = errors.New("nickname is already in use")

// GORM 및 일반 오류 체크 헬퍼
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound) || errors.Is(err, gorm.ErrRecordNotFound)
}

func IsAlreadyExists(err error) bool {
	return errors.Is(err, ErrAlreadyExists)
}
