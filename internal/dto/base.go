package dto

type Period struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type Pagination struct {
	Page       int64 `json:"page"`
	Limit      int64 `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}
