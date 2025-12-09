package dto

import (
	"github.com/google/uuid"
	"github.com/rainbow96bear/planet_utils/pb"
)

// 회원가입 요청 DTO
type CreateProfileRequest struct {
	UserID       uuid.UUID `json:"user_id"`
	Nickname     string    `json:"nickname"`
	Bio          string    `json:"bio"`
	ProfileImage string    `json:"profile_image"`
	Theme        string    `json:"theme"` // JSON 문자열 또는 미리 정의된 테마 값
}

// 회원가입 응답 DTO
type ProfileResponse struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Nickname     string    `json:"nickname"`
	Bio          string    `json:"bio"`
	ProfileImage string    `json:"profile_image"`
	Theme        string    `json:"theme"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
}

type ProfileUpdate struct {
	UserID       uuid.UUID `json:"user_id"`                 // DB의 user_id와 매핑
	Nickname     *string   `json:"nickname,omitempty"`      // nil이면 업데이트하지 않음
	Bio          *string   `json:"bio,omitempty"`           // nil이면 업데이트하지 않음
	ProfileImage *string   `json:"profile_image,omitempty"` // nil이면 업데이트하지 않음
	Theme        *string   `json:"Theme,omitempty"`         // nil이면 업데이트하지 않음
}

type UserProfile struct {
	ID             string    `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	Nickname       string    `json:"nickname"`
	Bio            string    `json:"bio,omitempty"`
	ProfileImage   string    `json:"profile_image,omitempty"`
	Theme          string    `json:"theme"`
	FollowerCount  int32     `json:"follower_count"`
	FollowingCount int32     `json:"following_count"`
}

func FromGrpcCreateUserRequest(req *pb.CreateUserRequest) (CreateProfileRequest, error) {
	profileImage := ""
	if req.ProfileImage != nil {
		profileImage = *req.ProfileImage
	}

	bio := ""
	if req.Bio != nil {
		bio = *req.Bio
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return CreateProfileRequest{}, err
	}

	return CreateProfileRequest{
		UserID:       userID,
		Nickname:     req.Nickname,
		Bio:          bio,
		ProfileImage: profileImage,
		Theme:        "light", // 기본값 문자열
	}, nil
}
