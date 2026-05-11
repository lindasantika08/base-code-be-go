package model

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`   
	Errors  interface{} `json:"errors,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type PaginationQuery struct {
	Page  int    `form:"page"   binding:"omitempty,min=1"`
	Limit int    `form:"limit"  binding:"omitempty,min=1,max=100"`
	Sort  string `form:"sort"`   
	Search string `form:"search"`
}

func (p *PaginationQuery) DefaultPagination() {
	if p.Page == 0 {
		p.Page = 1
	}
	if p.Limit == 0 {
		p.Limit = 10
	}
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

// ── HELPER CONSTRUCTORS ──────────────────────────────────────

func SuccessResponse(message string, data interface{}) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data}
}

func SuccessListResponse(message string, data interface{}, meta *Meta) APIResponse {
	return APIResponse{Success: true, Message: message, Data: data, Meta: meta}
}

func ErrorResponse(message string, errors interface{}) APIResponse {
	return APIResponse{Success: false, Message: message, Errors: errors}
}
