package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/rasadov/subscription-manager/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscription(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, id int) error
	ListSubscriptions(ctx context.Context,
		page, elements int,
		userID *string, serviceName *string,
		startDateFrom *time.Time, startDateTo *time.Time,
		endDateFrom *time.Time, endDateTo *time.Time,
		sortBy *string, sortOrder *string) (subscriptions []*models.Subscription, total int64, err error)
	CalculateTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate *time.Time) (int64, error)
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

	if res.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &subscription, nil
}

func (s *subscriptionRepository) UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error {
	return s.db.WithContext(ctx).Save(subscription).Error
}

func (s *subscriptionRepository) DeleteSubscription(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&models.Subscription{}, id).Error
}

func (s *subscriptionRepository) ListSubscriptions(ctx context.Context,
	page, elements int, userID *string, serviceName *string,
	startDateFrom *time.Time, startDateTo *time.Time,
	endDateFrom *time.Time, endDateTo *time.Time,
	sortBy *string, sortOrder *string) (subscriptions []*models.Subscription, total int64, err error) {
	db := s.db.WithContext(ctx)

	if userID != nil {
		db = db.Where("user_id = ?", *userID)
	}

	if startDateFrom != nil {
		db = db.Where("start_date >= ?", startDateFrom)
	}

	if startDateTo != nil {
		db = db.Where("start_date <= ?", startDateTo)
	}

	if endDateFrom != nil {
		db = db.Where("end_date >= ?", endDateFrom)
	}

	if endDateTo != nil {
		db = db.Where("end_date <= ?", endDateTo)
	}

	if serviceName != nil {
		db = db.Where("service_name = ?", *serviceName)
	}

	// Count total records
	if err := db.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sorting
	if sortBy != nil && sortOrder != nil {
		db = db.Order(*sortBy + " " + *sortOrder)
	} else {
		db = db.Order("created_at desc")
	}

	// Pagination
	limit := 10
	if elements > 0 {
		limit = elements
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	db = db.Limit(limit).Offset(offset)

	if err := db.Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (s *subscriptionRepository) CalculateTotalCost(ctx context.Context, userID, serviceName string, startDate, endDate *time.Time) (int64, error) {
	var totalCost int64

	db := s.db.WithContext(ctx).Model(&models.Subscription{})

	if userID != "" {
		db = db.Where("user_id = ?", userID)
	}
	if serviceName != "" {
		db = db.Where("service_name = ?", serviceName)
	}
	if startDate != nil {
		db = db.Where("start_date >= ?", *startDate)
	}
	if endDate != nil {
		db = db.Where("end_date <= ?", *endDate)
	}

	var totalCostNull sql.NullInt64
	if err := db.Select("SUM(price) as total_cost").Scan(&totalCostNull).Error; err != nil {
		return 0, err
	}
	totalCost = int64(0)
	if totalCostNull.Valid {
		totalCost = totalCostNull.Int64
	}

	return totalCost, nil
}
