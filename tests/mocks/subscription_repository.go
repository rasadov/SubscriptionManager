package mocks

import (
	"context"
	"time"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/repository"
)

type mockSubscriptionRepository struct {
}

var (
	subscriptions = []*models.Subscription{
		getMockSubscription(1),
		getMockSubscription(2),
		getMockSubscription(3),
		getMockSubscription(4),
		getMockSubscription(5),
	}
)

func NewMockSubscriptionRepository() repository.SubscriptionRepository {
	return &mockSubscriptionRepository{}
}

func getMockSubscription(id uint) *models.Subscription {
	return &models.Subscription{
		ID:          id,
		ServiceName: "test",
		Price:       10,
		UserID:      "test",
		StartDate:   "2025-01-01",
		EndDate:     "2025-01-02",
	}
}

func (m *mockSubscriptionRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	subscription.ID = 1
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	return nil
}

func (m *mockSubscriptionRepository) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	subscription := subscriptions[id-1]
	return subscription, nil
}

func (m *mockSubscriptionRepository) UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error {
	subscription.UpdatedAt = time.Now()
	return nil
}

func (m *mockSubscriptionRepository) DeleteSubscription(ctx context.Context, id int) error {
	return nil
}

func (m *mockSubscriptionRepository) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (subscriptions []*models.Subscription, total int64, err error) {
	return subscriptions, 5, nil
}

func (m *mockSubscriptionRepository) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	startPeriod, endPeriod := "2025-01-01 12:02", "2025-01-02 12:02"
	return &dto.TotalCostResponse{
		TotalCost: 10,
		Period: &dto.Period{
			StartDate: &startPeriod,
			EndDate:   &endPeriod,
		},
	}, nil
}
