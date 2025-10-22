package dto

type ProfileInfo struct {
	UserUuid     string `json:"user_uuid" db:"user_uuid"`
	Nickname     string `json:"nickname" db:"nickname"`
	ProfileImage string `json:"profile_image" db:"profile_image"`
	Bio          string `json:"bio,omitempty" db:"bio"`
	Email        string `json:"email,omitempty" db:"email"`
	IsOwner      bool   `json:"is_owner"`
}
