package repository

import (
	"context"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscription(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, id int) error
	ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) ([]models.Subscription, error)
	CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error)
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepositiry(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (s *subscriptionRepository) CreateSubscription(ctx context.Context, subscription *models.Subscription) error {
	return s.db.WithContext(ctx).Create(subscription).Error
}

func (s *subscriptionRepository) GetSubscription(ctx context.Context, id int) (*models.Subscription, error) {
	var subscription models.Subscription

	res := s.db.WithContext(ctx).Find(&subscription, id)
	if res.Error != nil {
		return nil, res.Error
	}

	return &subscription, nil
}

func (s *subscriptionRepository) UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error {
	return s.db.WithContext(ctx).Save(subscription).Error
}

func (s *subscriptionRepository) DeleteSubscription(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&models.Subscription{}, id).Error
}

func (s *subscriptionRepository) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) ([]models.Subscription, error) {
	return nil, nil
}

func (s *subscriptionRepository) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	return nil, nil
}
