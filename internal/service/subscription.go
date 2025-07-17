package service

import (
	"context"
	"time"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	GetSubscription(ctx context.Context, id int) (*dto.SubscriptionResponse, error)
	UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	DeleteSubscription(ctx context.Context, id int) error
	ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, error)
	CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{repo: repo}
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	startDate, err := time.Parse("01-2006", req.StartDate)
	if err != nil {
		return nil, err
	}
	var endDatePtr *time.Time
	if req.EndDate != "" {
		endDate, err := time.Parse("01-2006", req.EndDate)
		if err != nil {
			return nil, err
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
		return nil, err
	}

	return dto.NewSubscriptionResponse(&subscription), nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*dto.SubscriptionResponse, error) {
	subscription, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	return dto.NewSubscriptionResponse(subscription), nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	subscription, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
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
			return nil, err
		}
		subscription.StartDate = startDate
	}

	if req.EndDate != nil {
		endDate, err := time.Parse("01-2006", *req.EndDate)
		if err != nil {
			return nil, err
		}
		subscription.EndDate = &endDate
	}

	if err := s.repo.UpdateSubscription(ctx, id, subscription); err != nil {
		return nil, err
	}

	return dto.NewSubscriptionResponse(subscription), nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.DeleteSubscription(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}
	subscriptions, total, err := s.repo.ListSubscriptions(ctx, query)
	if err != nil {
		return nil, err
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
			Total:      total,
			TotalPages: (total + query.Limit - 1) / query.Limit,
		},
	}, nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	return s.repo.CalculateTotalCost(ctx, query)
}
