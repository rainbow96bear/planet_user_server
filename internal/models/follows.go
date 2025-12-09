package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	FollowerID  uuid.UUID `gorm:"type:uuid;not null;index"` // follower
	FollowingID uuid.UUID `gorm:"type:uuid;not null;index"` // following
	CreatedAt   time.Time `gorm:"autoCreateTime"`

	// 관계 정의 (optional)
	Follower  Profile `gorm:"foreignKey:FollowerID;references:UserID;constraint:OnDelete:CASCADE"`
	Following Profile `gorm:"foreignKey:FollowingID;references:UserID;constraint:OnDelete:CASCADE"`
}

func (Follow) TableName() string {
	return "follows"
}
