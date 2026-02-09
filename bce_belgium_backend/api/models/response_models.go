package models

type APIResponse struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Count    int   `json:"count,omitempty"`
	Total    int   `json:"total,omitempty"`
	Limit    int   `json:"limit,omitempty"`
	Offset   int   `json:"offset,omitempty"`
	Page     int   `json:"page,omitempty"`
	Pages    int   `json:"pages,omitempty"`
	Duration int64 `json:"duration_ms,omitempty"`
}

type PaginatedResponse struct {
	Data       any  `json:"data"`
	Pagination Meta `json:"pagination"`
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

func Paginated(data any, meta Meta) PaginatedResponse {
	return PaginatedResponse{
		Data:       data,
		Pagination: meta,
	}
}
