package dto

import "mime/multipart"

type TodoItem struct {
	Text      string `json:"text"`
	Completed bool   `json:"completed"`
}

type CalendarInfo struct {
	EventID     int64      `json:"event_id"`
	UserUUID    string     `json:"user_uuid"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Emoji       string     `json:"emoji"`
	StartAt     string     `json:"start_at"`
	EndAt       string     `json:"end_at"`
	Visibility  string     `json:"visibility"` // public | friends | private
	ImageURL    string     `json:"image_url"`  // 현재는 빈 값, 추후 사용
	Todos       []TodoItem `json:"todos"`
}

type CalendarCreateRequest struct {
	Title       string `form:"title" binding:"required"`
	Emoji       string `form:"emoji"`
	StartDate   string `form:"startDate" binding:"required"`
	EndDate     string `form:"endDate" binding:"required"`
	Description string `form:"description"`
	Visibility  string `form:"visibility" binding:"required"`
	Todos       string `form:"todos"` // JSON 문자열 형태
	// 추후 이미지 업로드 활성화 시 주석 해제
	// Image *multipart.FileHeader `form:"image"`
}

// 추후 이미지 업로드 지원을 위한 헬퍼 함수 (현재는 사용 안함)
func (r *CalendarCreateRequest) HasImage() bool {
	// return r.Image != nil && r.Image.Size > 0
	return false // MVP에서는 항상 false
}

type CalendarUpdateRequest struct {
	Title       string `form:"title"`
	Description string `form:"description"`
	Emoji       string `form:"emoji"`
	StartDate   string `form:"startDate"`
	EndDate     string `form:"endDate"`
	Visibility  string `form:"visibility"`
	Todos       string `form:"todos"` // JSON 문자열 형태
	// 추후 이미지 업로드 활성화 시 주석 해제
	// Image       *multipart.FileHeader `form:"image"`
	// DeleteImage bool                  `form:"deleteImage"` // 기존 이미지 삭제 플래그
}

func ToCalendarModel(req *CalendarCreateRequest, userUUID string, todos []TodoItem) *CalendarInfo {
	return &CalendarInfo{
		UserUUID:    userUUID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     req.StartDate,
		EndAt:       req.EndDate,
		Visibility:  req.Visibility,
		ImageURL:    "", // MVP에서는 이미지 등록 안함
		Todos:       todos,
	}
}

func ToCalendarUpdateModel(req *CalendarUpdateRequest, userUUID string, eventID int64, todos []TodoItem, imageURL string) *CalendarInfo {
	return &CalendarInfo{
		EventID:     eventID,
		UserUUID:    userUUID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     req.StartDate,
		EndAt:       req.EndDate,
		Visibility:  req.Visibility,
		ImageURL:    imageURL, // 기존 이미지 URL 유지 또는 새 URL
		Todos:       todos,
	}
}

// 추후 이미지 업로드 서비스 구현 시 사용할 인터페이스
type ImageUploader interface {
	Upload(file *multipart.FileHeader) (string, error)
	Delete(url string) error
}
