package dto

import "github.com/rasadov/subscription-manager/internal/models"

type CreateSubscriptionRequest struct {
	ServiceName string  `json:"service_name" binding:"required"`
	Price       int64   `json:"price" binding:"required"`
	UserID      string  `json:"user_id" binding:"required,uuid"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     *string `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string `json:"service_name,omitempty"`
	Price       *int64  `json:"price,omitempty"`
	StartDate   *string `json:"start_date,omitempty"`
	EndDate     *string `json:"end_date,omitempty"`
}

type ListSubscriptionsQuery struct {
	UserID        *string `form:"user_id"`
	ServiceName   *string `form:"service_name"`
	Page          int64   `form:"page,default=1"`
	Limit         int64   `form:"limit,default=10"`
	StartDateFrom *string `form:"start_date_from"`
	EndDateFrom   *string `form:"end_date_from"`
	EndDateTo     *string `form:"end_date_to"`
	SortBy        *string `form:"sort_by"`
	SortOrder     *string `form:"sort_order"`
}

type TotalCostQuery struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	StartDate   *string `form:"start_date" binding:"required"`
	EndDate     *string `form:"end_date" binding:"required"`
}

type TotalCostResponse struct {
	TotalCost int64 `json:"total_cost"`
	Period    *Period
}

type ListSubscriptionsResponse struct {
	Data       []*models.Subscription `json:"data"`
	Pagination *Pagination            `json:"pagination"`
}
