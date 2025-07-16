package repository

import (
	"context"
	"database/sql"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, subscription *models.Subscription) error
	GetSubscription(ctx context.Context, id int) (*models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int, subscription *models.Subscription) error
	DeleteSubscription(ctx context.Context, id int) error
	ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (subscriptions []*models.Subscription, total int64, err error)
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

func (s *subscriptionRepository) ListSubscriptions(ctx context.Context, query dto.ListSubscriptionsQuery) (subscriptions []*models.Subscription, total int64, err error) {
	db := s.db.WithContext(ctx)

	if err := db.Model(&models.Subscription{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}
	if query.ServiceName != nil {
		db = db.Where("service_name = ?", *query.ServiceName)
	}
	if query.StartDateFrom != nil {
		db = db.Where("start_date >= ?", *query.StartDateFrom)
	}
	if query.EndDateFrom != nil {
		db = db.Where("end_date >= ?", *query.EndDateFrom)
	}
	if query.EndDateTo != nil {
		db = db.Where("end_date <= ?", *query.EndDateTo)
	}

	// Sorting
	sortBy := "created_at"
	if query.SortBy != nil && *query.SortBy != "" {
		sortBy = *query.SortBy
	}
	sortOrder := "desc"
	if query.SortOrder != nil && (*query.SortOrder == "asc" || *query.SortOrder == "desc") {
		sortOrder = *query.SortOrder
	}
	db = db.Order(sortBy + " " + sortOrder)

	// Pagination
	limit := 10
	if query.Limit > 0 {
		limit = int(query.Limit)
	}
	page := 1
	if query.Page > 0 {
		page = int(query.Page)
	}
	offset := (page - 1) * limit
	db = db.Limit(limit).Offset(offset)

	if err := db.Find(&subscriptions).Error; err != nil {
		return nil, 0, err
	}

	return subscriptions, total, nil
}

func (s *subscriptionRepository) CalculateTotalCost(ctx context.Context, query dto.TotalCostQuery) (*dto.TotalCostResponse, error) {
	var totalCost int64

	db := s.db.WithContext(ctx).Model(&models.Subscription{})

	if query.UserID != nil {
		db = db.Where("user_id = ?", *query.UserID)
	}
	if query.ServiceName != nil {
		db = db.Where("service_name = ?", *query.ServiceName)
	}
	if query.StartDate != nil {
		db = db.Where("start_date >= ?", *query.StartDate)
	}
	if query.EndDate != nil {
		db = db.Where("end_date <= ?", *query.EndDate)
	}

	var totalCostNull sql.NullInt64
	if err := db.Select("SUM(price) as total_cost").Scan(&totalCostNull).Error; err != nil {
		return nil, err
	}
	totalCost = int64(0)
	if totalCostNull.Valid {
		totalCost = totalCostNull.Int64
	}

	return &dto.TotalCostResponse{
		TotalCost: totalCost,
		Period: &dto.Period{
			StartDate: query.StartDate,
			EndDate:   query.EndDate,
		},
	}, nil
}
