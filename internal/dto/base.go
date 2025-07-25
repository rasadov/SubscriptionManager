package dto

import (
	"fmt"
	"time"
)

type Period struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
}

type Pagination struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

type MonthYear time.Time

func (m MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(m)
	if t.IsZero() {
		return []byte("null"), nil
	}
	s := fmt.Sprintf("\"%02d-%04d\"", t.Month(), t.Year())
	return []byte(s), nil
}
