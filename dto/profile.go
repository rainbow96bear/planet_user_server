package dto

import "github.com/google/uuid"

// -------------------------------
// 클라이언트 요청 DTO
// -------------------------------

// 프로필 업데이트 요청
type ProfileUpdateRequest struct {
	Nickname     string `json:"nickname" binding:"required"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Theme        string `json:"theme"` // optional
}

// 테마만 수정하는 요청
type ThemeUpdateRequest struct {
	Theme string `json:"theme" binding:"required"`
}

// -------------------------------
// 내부/서비스 공용 DTO (도메인 모델과 1:1 대응)
// -------------------------------

type ProfileInfo struct {
	UserID       uuid.UUID `json:"user_id"`
	Nickname     string    `json:"nickname"`
	Bio          string    `json:"bio"`
	ProfileImage string    `json:"profile_image"`
	Theme        string    `json:"theme"`

	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`
}

// 서비스/DB 레이어에서 업데이트 처리용
type ProfileUpdate struct {
	UserID       uuid.UUID `json:"user_id"`
	Nickname     *string   `json:"nickname"`      // nil이면 업데이트 안 함
	Bio          *string   `json:"bio"`           // nil이면 업데이트 안 함
	ProfileImage *string   `json:"profile_image"` // nil이면 업데이트 안 함
	Theme        *string   `json:"theme"`         // nil이면 업데이트 안 함
}

// -------------------------------
// 클라이언트 응답 DTO
// -------------------------------

// 단일 프로필 조회 응답
type ProfileResponse struct {
	UserID       string `json:"user_id"`
	Nickname     string `json:"nickname"`
	Bio          string `json:"bio"`
	ProfileImage string `json:"profile_image"`
	Theme        string `json:"theme"`

	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`
}

// 팔로워/팔로잉 수 응답
type FollowCountDTO struct {
	FollowerCount  int `json:"follower_count"`
	FollowingCount int `json:"following_count"`
}

// -------------------------------
// 변환 함수
// -------------------------------

// 요청 DTO → 내부 DTO 변환 (업데이트용)
func ToProfileUpdateModel(req *ProfileUpdateRequest, UserID uuid.UUID) *ProfileUpdate {
	update := &ProfileUpdate{
		UserID: UserID,
	}

	// 필드별 nil 처리
	if req.Nickname != "" {
		update.Nickname = &req.Nickname
	}
	if req.Bio != "" {
		update.Bio = &req.Bio
	}
	if req.ProfileImage != "" {
		update.ProfileImage = &req.ProfileImage
	}
	if req.Theme != "" {
		update.Theme = &req.Theme
	}

	return update
}

// 내부 DTO → 클라이언트 응답 DTO 변환
func ToProfileResponse(info *ProfileInfo) *ProfileResponse {
	if info == nil {
		return &ProfileResponse{}
	}

	return &ProfileResponse{
		UserID:         info.UserID.String(),
		Nickname:       info.Nickname,
		Bio:            info.Bio,
		ProfileImage:   info.ProfileImage,
		Theme:          info.Theme,
		FollowerCount:  info.FollowerCount,
		FollowingCount: info.FollowingCount,
	}
}
