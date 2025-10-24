package dto

// 클라이언트 요청용 DTO (업데이트)
type ProfileUpdateRequest struct {
	Nickname     string `json:"nickname,omitempty"`
	Bio          string `json:"bio,omitempty"`
	Email        string `json:"email,omitempty"`
	ProfileImage string `json:"profile_image,omitempty"`
}

// 내부/서비스용 DTO (service, repository 로직 공통)
type ProfileInfo struct {
	UserUuid     string
	Nickname     string
	Bio          string
	Email        string
	ProfileImage string
}

// 조회용 응답 DTO
type ProfileResponse struct {
	Nickname     string `json:"nickname"`
	Bio          string `json:"bio"`
	Email        string `json:"email,omitempty"`
	ProfileImage string `json:"profile_image"`
}

func ToProfileInfo(req *ProfileUpdateRequest, userUuid string) *ProfileInfo {
	return &ProfileInfo{
		UserUuid:     userUuid,
		Nickname:     req.Nickname,
		Bio:          req.Bio,
		Email:        req.Email,
		ProfileImage: req.ProfileImage,
	}
}

func ToProfileResponse(req *ProfileInfo) *ProfileResponse {
	return &ProfileResponse{
		Nickname:     req.Nickname,
		Bio:          req.Bio,
		Email:        req.Email,
		ProfileImage: req.ProfileImage,
	}
}
