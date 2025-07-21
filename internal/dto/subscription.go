package dto

import (
	"time"

	"github.com/rasadov/subscription-manager/internal/models"
)

type CreateSubscriptionRequest struct {
	ServiceName string `json:"service_name" binding:"required"`
	Price       int64  `json:"price" binding:"required"`
	UserID      string `json:"user_id" binding:"required,uuid"`
	StartDate   string `json:"start_date" binding:"required"`
	EndDate     string `json:"end_date,omitempty"`
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
	Page          int     `form:"page,default=1"`
	Limit         int     `form:"limit,default=10"`
	StartDateFrom *string `form:"start_date_from"`
	StartDateTo   *string `form:"start_date_to"`
	EndDateFrom   *string `form:"end_date_from"`
	EndDateTo     *string `form:"end_date_to"`
	SortBy        *string `form:"sort_by"`
	SortOrder     *string `form:"sort_order"`
}

type SubscriptionResponse struct {
	ID          uint       `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int64      `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   MonthYear  `json:"start_date"`
	EndDate     *MonthYear `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
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
	Data       []*SubscriptionResponse `json:"data"`
	Pagination *Pagination             `json:"pagination"`
}

func NewSubscriptionResponse(subscription *models.Subscription) *SubscriptionResponse {
	var endDate *MonthYear
	if subscription.EndDate != nil {
		endDateVal := MonthYear(*subscription.EndDate)
		endDate = &endDateVal
	}

	return &SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   MonthYear(subscription.StartDate),
		EndDate:     endDate,
		CreatedAt:   subscription.CreatedAt,
		UpdatedAt:   subscription.UpdatedAt,
	}
}
