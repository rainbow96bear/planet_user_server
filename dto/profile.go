package dto

type profileInfo struct {
	Nickname      string    `json:"nickname" db:"nickname"`               // 닉네임
	ProfileImage  string    `json:"profile_image" db:"profile_image"`     // 프로필 이미지 URL
	Bio           string    `json:"bio,omitempty" db:"bio"`               // 자기소개
	Email         string    `json:"email,omitempty" db:"email"`           // 이메일 (NULL 가능)
}
