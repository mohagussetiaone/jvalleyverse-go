package dto

// APIResponse standard response wrapper for all API responses
type APIResponse struct {
	Data       interface{} `json:"data"`
	Message    string      `json:"message,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

// Pagination holds pagination metadata
type Pagination struct {
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
	Total int64 `json:"total"`
}
