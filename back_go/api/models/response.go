package models

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Count  int `json:"count,omitempty"`
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`
}

func Success(data any) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
	}
}

func SuccessWithMeta(data any, meta Meta) APIResponse {
	return APIResponse{
		Success: true,
		Data:    data,
		Meta:    &meta,
	}
}

func Error(message string) APIResponse {
	return APIResponse{
		Success: false,
		Error:   message,
	}
}
