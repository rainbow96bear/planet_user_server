package dto

type UpdateProfileRequest struct {
	Nickname string `form:"nickname" binding:"required"`
	Bio      string `form:"bio"`
	Email    string `form:"email"`
}
