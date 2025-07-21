package service

import (
	"context"
	"errors"
	"time"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/repository"
	"github.com/rasadov/subscription-manager/pkg/exceptions"
	"gorm.io/gorm"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, exceptions.HTTPError)
	GetSubscription(ctx context.Context, id int) (*dto.SubscriptionResponse, exceptions.HTTPError)
	UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, exceptions.HTTPError)
	DeleteSubscription(ctx context.Context, id int) exceptions.HTTPError
	ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, exceptions.HTTPError)
	CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, exceptions.HTTPError)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, exceptions.HTTPError) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, exceptions.NewBadRequest(err.Error())
	}
	var endDatePtr *time.Time
	if req.EndDate != "" {
		endDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		endDatePtr = &endDate
	}
	subscription := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   startDate,
		EndDate:     endDatePtr,
	}

	if err := s.repo.CreateSubscription(ctx, &subscription); err != nil {
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	return dto.NewSubscriptionResponse(&subscription), nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*dto.SubscriptionResponse, exceptions.HTTPError) {
	subscription, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFound(err.Error())
		}
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	return dto.NewSubscriptionResponse(subscription), nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, exceptions.HTTPError) {
	subscription, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.NewNotFound(err.Error())
		}
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	if req.ServiceName != nil {
		subscription.ServiceName = *req.ServiceName
	}

	if req.Price != nil {
		subscription.Price = *req.Price
	}

	if req.StartDate != nil {
		startDate, err := time.Parse("01-2006", *req.StartDate)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		subscription.StartDate = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		subscription.EndDate = &endDate
	}

	if err := s.repo.UpdateSubscription(ctx, id, subscription); err != nil {
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	return dto.NewSubscriptionResponse(subscription), nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) exceptions.HTTPError {
	err := s.repo.DeleteSubscription(ctx, id)
	if err != nil {
		return exceptions.NewInternalServerError(err.Error())
	}
	return nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, exceptions.HTTPError) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	sortBy := "created_at"
	if query.SortBy != nil && *query.SortBy != "" {
		sortBy = *query.SortBy
	}

	sortOrder := "desc"
	if query.SortOrder != nil && (*query.SortOrder == "asc" || *query.SortOrder == "desc") {
		sortOrder = *query.SortOrder
	}

	var startDateFrom *time.Time
	var startDateTo *time.Time
	var endDateFrom *time.Time
	var endDateTo *time.Time
	var err error

	if query.StartDateFrom != nil {
		startDateFromParsed, err := time.Parse("01-2006", *query.StartDateFrom)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		startDateFrom = &startDateFromParsed
	}

	if query.StartDateTo != nil {
		startDateToParsed, err := time.Parse("01-2006", *query.StartDateTo)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		startDateTo = &startDateToParsed
	}

	if query.EndDateFrom != nil {
		endDateFromParsed, err := time.Parse("01-2006", *query.EndDateFrom)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		endDateFrom = &endDateFromParsed
	}

	if query.EndDateTo != nil {
		endDateToParsed, err := time.Parse("01-2006", *query.EndDateTo)
		if err != nil {
			return nil, exceptions.NewBadRequest(err.Error())
		}
		endDateTo = &endDateToParsed
	}

	subscriptions, total, err := s.repo.ListSubscriptions(ctx, int(query.Page), int(query.Limit), query.UserID, query.ServiceName,
		startDateFrom, startDateTo, endDateFrom, endDateTo, &sortBy, &sortOrder)
	if err != nil {
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	var subscriptionResponses []*dto.SubscriptionResponse
	for _, subscription := range subscriptions {
		subscriptionResponses = append(subscriptionResponses, dto.NewSubscriptionResponse(subscription))
	}

	return &dto.ListSubscriptionsResponse{
		Data: subscriptionResponses,
		Pagination: &dto.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      int(total),
			TotalPages: (int(total) + query.Limit - 1) / query.Limit,
		},
	}, nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, exceptions.HTTPError) {
	startDateParsed, err := time.Parse("01-2006", *query.StartDate)
	if err != nil {
		return nil, exceptions.NewBadRequest(err.Error())
	}
	endDateParsed, err := time.Parse("01-2006", *query.EndDate)
	if err != nil {
		return nil, exceptions.NewBadRequest(err.Error())
	}

	totalCost, err := s.repo.CalculateTotalCost(ctx, *query.UserID, *query.ServiceName, &startDateParsed, &endDateParsed)
	if err != nil {
		return nil, exceptions.NewInternalServerError(err.Error())
	}

	return &dto.TotalCostResponse{
		TotalCost: totalCost,
		Period: &dto.Period{
			StartDate: query.StartDate,
			EndDate:   query.EndDate,
		},
	}, nil
}
