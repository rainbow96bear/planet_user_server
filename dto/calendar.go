package dto

type CalendarInfo struct {
	EventID     int64  `json:"event_id"`
	UserUUID    string `json:"user_uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	Visibility  string `json:"visibility"` // public | friends | private
	ImageURL    string `json:"image_url"`
}

type CalendarCreateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	Visibility  string `json:"visibility"`
	ImageURL    string `json:"image_url"`
}

func ToCalendarModel(req *CalendarCreateRequest, userUUID string) *CalendarInfo {
	return &CalendarInfo{
		UserUUID:    userUUID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Visibility:  req.Visibility,
		ImageURL:    req.ImageURL,
	}
}

type CalendarUpdateRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Emoji       string `json:"emoji"`
	StartAt     string `json:"start_at"`
	EndAt       string `json:"end_at"`
	Visibility  string `json:"visibility"`
	ImageURL    string `json:"image_url"`
}

func ToCalendarUpdateModel(req *CalendarUpdateRequest, userUUID string, eventID int64) *CalendarInfo {
	return &CalendarInfo{
		EventID:     eventID,
		UserUUID:    userUUID,
		Title:       req.Title,
		Description: req.Description,
		Emoji:       req.Emoji,
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		Visibility:  req.Visibility,
		ImageURL:    req.ImageURL,
	}
}
