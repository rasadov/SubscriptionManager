package service

import (
	"context"

	"github.com/rasadov/subscription-manager/internal/dto"
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
	return nil, nil
}

func (s *subscriptionService) GetSubscription(ctx context.Context, id int) (*dto.SubscriptionResponse, error) {
	return nil, nil
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int, req dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	return nil, nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, id int) error {
	return nil
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (*dto.ListSubscriptionsResponse, error) {
	return nil, nil
}

func (s *subscriptionService) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	return nil, nil
}
