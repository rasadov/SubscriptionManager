package mocks

import (
	"context"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockSubscriptionRepository struct {
	mock.Mock
}

func (m *MockSubscriptionRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	args := m.Called(ctx, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockSubscriptionRepository) UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error {
	args := m.Called(ctx, id, subscription)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) DeleteSubscription(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSubscriptionRepository) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) ([]*models.Subscription, int64, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]*models.Subscription), args.Get(1).(int64), args.Error(2)
}

func (m *MockSubscriptionRepository) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TotalCostResponse), args.Error(1)
}
