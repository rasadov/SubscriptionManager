package tests

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	repoTestServiceName = "repo-test"
	repoTestPrice       = int64(100)
	repoTestUserID      = "repo-user"
	repoTestStartDate   = "2025-01-01"
	repoTestEndDate     = "2025-01-31"
)

func createTestSubscription(t *testing.T) *models.Subscription {
	sub := &models.Subscription{
		ServiceName: repoTestServiceName,
		Price:       repoTestPrice,
		UserID:      repoTestUserID,
		StartDate:   repoTestStartDate,
		EndDate:     repoTestEndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := testRepository.CreateSubscription(context.Background(), sub)
	assert.NoError(t, err)
	return sub
}

func TestCreateAndGetSubscriptionRepository(t *testing.T) {
	SetupRepo(t)
	sub := createTestSubscription(t)

	got, err := testRepository.GetSubscription(context.Background(), int(sub.ID))
	assert.NoError(t, err)
	assert.Equal(t, sub.ServiceName, got.ServiceName)
	assert.Equal(t, sub.Price, got.Price)
	assert.Equal(t, sub.UserID, got.UserID)
}

func TestUpdateSubscriptionRepository(t *testing.T) {
	SetupRepo(t)
	sub := createTestSubscription(t)

	newServiceName := "updated-service"
	sub.ServiceName = newServiceName
	err := testRepository.UpdateSubscription(context.Background(), int(sub.ID), sub)
	assert.NoError(t, err)

	updated, err := testRepository.GetSubscription(context.Background(), int(sub.ID))
	assert.NoError(t, err)
	assert.Equal(t, newServiceName, updated.ServiceName)
}

func TestDeleteSubscriptionRepository(t *testing.T) {
	SetupRepo(t)
	sub := createTestSubscription(t)

	err := testRepository.DeleteSubscription(context.Background(), int(sub.ID))
	assert.NoError(t, err)

	sub, err = testRepository.GetSubscription(context.Background(), int(sub.ID))
	log.Println("err: ", err)
	log.Println("sub: ", sub)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	assert.Nil(t, sub)
}

func TestListSubscriptionsRepository(t *testing.T) {
	SetupRepo(t)
	_ = createTestSubscription(t)

	query := dto.ListSubscriptionsQuery{
		UserID: &repoTestUserID,
	}
	subs, total, err := testRepository.ListSubscriptions(context.Background(), query)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, subs, 1)
	assert.Equal(t, repoTestUserID, subs[0].UserID)
}

func TestCalculateTotalCostRepository(t *testing.T) {
	SetupRepo(t)
	_ = createTestSubscription(t)

	query := dto.TotalCostQuery{
		UserID:      &repoTestUserID,
		ServiceName: &repoTestServiceName,
		StartDate:   &repoTestStartDate,
		EndDate:     &repoTestEndDate,
	}
	resp, err := testRepository.CalculateTotalCost(context.Background(), query)
	assert.NoError(t, err)
	assert.Equal(t, repoTestPrice, resp.TotalCost)
}
