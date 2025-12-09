package models

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID           string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;unique" json:"user_id"`
	Nickname     string    `gorm:"size:50;not null;unique" json:"nickname"`
	Bio          string    `json:"bio,omitempty"`
	ProfileImage string    `json:"profile_image,omitempty"`
	Theme        string    `gorm:"size:20;not null;default:'light'" json:"theme"`

	// 팔로우 정보 (int32로 변경)
	FollowerCount  int32 `gorm:"not null;default:0" json:"follower_count"`
	FollowingCount int32 `gorm:"not null;default:0" json:"following_count"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (Profile) TableName() string {
	return "profiles"
}
