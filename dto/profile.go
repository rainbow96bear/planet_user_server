package dto

// 클라이언트 요청용 DTO (업데이트)
type ProfileUpdateRequest struct {
	Nickname     string `json:"nickname,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Email        string `json:"email,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}

// 내부/서비스/응답 공통용 DTO
type ProfileInfo struct {
	UserUUID       string `json:"uuid"`
	Nickname       string `json:"nickname"`
	ProfileImage   string `json:"profile_image"`
	Bio            string `json:"bio,omitempty"`
	Email          string `json:"email,omitempty"`
	FollowerCount  uint   `json:"followerCount"`
	FollowingCount uint   `json:"followingCount"`
}

// 조회용 응답 DTO
type ProfileResponse struct {
	Nickname       string `json:"nickname"`
	Bio            string `json:"bio,omitempty"`
	Email          string `json:"email,omitempty"`
	ProfileImage   string `json:"profile_image"`
	FollowerCount  uint   `json:"followerCount"`
	FollowingCount uint   `json:"followingCount"`
}

// 팔로워/팔로잉 수용 DTO
type FollowCountDTO struct {
	UserUuid       string `json:"user_uuid"`
	FollowerCount  uint   `json:"follower_count"`
	FollowingCount uint   `json:"following_count"`
}

type Theme struct {
	Theme string `json:"theme"`
}

// 요청 DTO → 내부 DTO 변환
func ToProfileInfo(req *ProfileUpdateRequest, userUuid string) *ProfileInfo {
	return &ProfileInfo{
		UserUUID:     userUuid,
		Nickname:     req.Nickname,
		Bio:          req.Bio,
		Email:        req.Email,
		ProfileImage: req.ProfileImage,
	}
}

// 내부 DTO → 응답 DTO 변환
func ToProfileResponse(info *ProfileInfo) *ProfileResponse {
	return &ProfileResponse{
		Nickname:       info.Nickname,
		Bio:            info.Bio,
		Email:          info.Email,
		ProfileImage:   info.ProfileImage,
		FollowerCount:  info.FollowerCount,
		FollowingCount: info.FollowingCount,
	}
}
