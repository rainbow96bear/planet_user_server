package planet_err

import "fmt"

// ErrorCode는 사용자 정의 에러 코드를 나타냅니다.
type ErrorCode string

// CodeError는 에러 코드와 메시지, HTTP 상태 코드, 원본 오류를 포함
type CodeError struct {
	Code    ErrorCode              `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"` // 추가 데이터
	Status  int                    `json:"-"`              // HTTP 상태 코드
	Err     error                  `json:"-"`              // 원본 오류
}

// NewCodeError 생성자
func NewCodeError(code ErrorCode, msg string, status int, originalErr error) *CodeError {
	return &CodeError{
		Code:    code,
		Message: msg,
		Status:  status,
		Err:     originalErr,
	}
}

// error 인터페이스 구현
func (e *CodeError) Error() string {
	return fmt.Sprintf("[%s] %s (Status: %d, Original: %v)", e.Code, e.Message, e.Status, e.Err)
}

// WithData: 추가 정보를 포함한 CodeError 반환
func (e *CodeError) WithData(data map[string]interface{}) *CodeError {
	e.Data = data
	return e
}
