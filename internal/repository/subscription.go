package repository

import "gorm.io/gorm"

type SubscriptionRepository interface {
}

type subscriptionRepository struct {
	db *gorm.DB
}

func NewSubscriptionRepositiry(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{db: db}
}
