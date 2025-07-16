package service

import (
	"context"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/repository"
)

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*models.Subscription, error)
	GetSubscription(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*models.Subscription, error)
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

func (s *subscriptionService) CreateSubscription(ctx context.Context, req dto.CreateSubscriptionRequest) (*models.Subscription, error) {
	subscription := models.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
	}

	if err := s.repo.CreateSubscription(ctx, &subscription); err != nil {
		return nil, err
	}

	return &subscription, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	subscription, err := s.repo.GetSubscription(ctx, id)
	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*models.Subscription, error) {
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
		subscription.StartDate = *req.StartDate
	}

	if req.EndDate != nil {
		subscription.EndDate = *req.EndDate
	}

	if err := s.repo.UpdateSubscription(ctx, id, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return s.repo.DeleteSubscription(ctx, id)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, error) {
	subscriptions, total, err := s.repo.ListSubscriptions(ctx, query)
	if err != nil {
		return nil, err
	}

	return &dto.ListSubscriptionsResponse{
		Data: subscriptions,
		Pagination: &dto.Pagination{
			Page:       query.Page,
			Limit:      query.Limit,
			Total:      total,
			TotalPages: (total + query.Limit - 1) / query.Limit,
		},
	}, nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	return nil, nil
}
