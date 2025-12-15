package models

// Response 统一响应结构
type Response struct {
	Success bool        `json:"success"`         // 是否成功
	Message string      `json:"message"`         // 响应消息
	Data    interface{} `json:"data,omitempty"`  // 响应数据（可选）
	Error   string      `json:"error,omitempty"` // 错误信息（可选）
}

// SuccessResponse 成功响应
func SuccessResponse(message string, data interface{}) Response {
	return Response{
		Success: true,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(message string, err error) Response {
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	return Response{
		Success: false,
		Message: message,
		Error:   errMsg,
	}
}

