package tests

import (
	"context"
	"testing"
	"time"

	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func createTestSubscription(serviceName, userID, startDate, endDate string, price int64) *models.Subscription {
	startDateParsed, _ := time.Parse("01-2006", startDate)
	endDateParsed, _ := time.Parse("01-2006", endDate)
	return &models.Subscription{
		ServiceName: serviceName,
		Price:       price,
		UserID:      userID,
		StartDate:   startDateParsed,
		EndDate:     &endDateParsed,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

func seedTestSubscriptions(t *testing.T) []*models.Subscription {
	subscriptions := []*models.Subscription{
		createTestSubscription("Netflix", "user-1", "01-2025", "12-2025", 1500),
		createTestSubscription("Spotify", "user-1", "02-2025", "11-2025", 1200),
		createTestSubscription("Netflix", "user-2", "03-2025", "10-2025", 1500),
		createTestSubscription("Apple Music", "user-2", "01-2025", "", 990),
		createTestSubscription("Disney+", "user-1", "06-2025", "12-2025", 899),
	}

	for _, sub := range subscriptions {
		err := testRepository.CreateSubscription(context.Background(), sub)
		require.NoError(t, err)
	}

	return subscriptions
}

func TestCreateSubscription_Success(t *testing.T) {
	SetupRepo(t)

	sub := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)

	err := testRepository.CreateSubscription(context.Background(), sub)

	assert.NoError(t, err)
	assert.NotZero(t, sub.ID)
	assert.False(t, sub.CreatedAt.IsZero())
	assert.False(t, sub.UpdatedAt.IsZero())
}

func TestGetSubscription_Success(t *testing.T) {
	SetupRepo(t)

	original := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)
	err := testRepository.CreateSubscription(context.Background(), original)
	require.NoError(t, err)

	retrieved, err := testRepository.GetSubscription(context.Background(), int(original.ID))

	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, original.ID, retrieved.ID)
	assert.Equal(t, original.ServiceName, retrieved.ServiceName)
	assert.Equal(t, original.Price, retrieved.Price)
	assert.Equal(t, original.UserID, retrieved.UserID)
	assert.Equal(t, original.StartDate, retrieved.StartDate)
	assert.Equal(t, original.EndDate, retrieved.EndDate)
}

func TestGetSubscription_NotFound(t *testing.T) {
	SetupRepo(t)

	retrieved, err := testRepository.GetSubscription(context.Background(), 999999)

	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestGetSubscription_InvalidID(t *testing.T) {
	SetupRepo(t)

	retrieved, err := testRepository.GetSubscription(context.Background(), -1)

	assert.Error(t, err)
	assert.Nil(t, retrieved)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestUpdateSubscription_Success(t *testing.T) {
	SetupRepo(t)

	original := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)
	err := testRepository.CreateSubscription(context.Background(), original)
	require.NoError(t, err)

	original.ServiceName = "Netflix Premium"
	original.Price = 2000
	err = testRepository.UpdateSubscription(context.Background(), int(original.ID), original)
	assert.NoError(t, err)

	updated, err := testRepository.GetSubscription(context.Background(), int(original.ID))
	require.NoError(t, err)
	assert.Equal(t, "Netflix Premium", updated.ServiceName)
	assert.Equal(t, int64(2000), updated.Price)
	assert.True(t, updated.UpdatedAt.After(updated.CreatedAt))
}

func TestUpdateSubscription_NotFound(t *testing.T) {
	SetupRepo(t)

	nonExistentSub := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)
	nonExistentSub.ID = 999999

	err := testRepository.UpdateSubscription(context.Background(), 999999, nonExistentSub)

	assert.NoError(t, err)
}

func TestDeleteSubscription_Success(t *testing.T) {
	SetupRepo(t)

	sub := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)
	err := testRepository.CreateSubscription(context.Background(), sub)
	require.NoError(t, err)

	err = testRepository.DeleteSubscription(context.Background(), int(sub.ID))
	assert.NoError(t, err)

	deleted, err := testRepository.GetSubscription(context.Background(), int(sub.ID))
	assert.Error(t, err)
	assert.Nil(t, deleted)
	assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
}

func TestDeleteSubscription_NotFound(t *testing.T) {
	SetupRepo(t)

	err := testRepository.DeleteSubscription(context.Background(), 999999)

	assert.NoError(t, err)
}

func TestListSubscriptions_NoFilters(t *testing.T) {
	SetupRepo(t)
	subs := seedTestSubscriptions(t)

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, nil, nil, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Len(t, result, len(subs))
	assert.Equal(t, int64(len(subs)), total)
}

func TestListSubscriptions_FilterByUserID(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	userID := "user-1"

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, &userID, nil, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(3), total)
	for _, sub := range result {
		assert.Equal(t, userID, sub.UserID)
	}
}

func TestListSubscriptions_FilterByServiceName(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	serviceName := "Netflix"

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, nil, &serviceName, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	for _, sub := range result {
		assert.Equal(t, serviceName, sub.ServiceName)
	}
}

func TestListSubscriptions_MultipleFilters(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	userID := "user-1"
	serviceName := "Netflix"

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, &userID, &serviceName, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, result, 1)
	assert.Equal(t, userID, result[0].UserID)
	assert.Equal(t, serviceName, result[0].ServiceName)
}

func TestListSubscriptions_DateFilters(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	startDateFrom := "02-2025"
	startDateFromParsed, _ := time.Parse("01-2006", startDateFrom)

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, nil, nil, &startDateFromParsed, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.True(t, total >= 1)
	for _, sub := range result {
		assert.True(t, sub.StartDate.After(startDateFromParsed) || sub.StartDate.Equal(startDateFromParsed))
	}
}

func TestListSubscriptions_Pagination(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	result1, total1, err := testRepository.ListSubscriptions(context.Background(), 1, 2, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, result1, 2)
	assert.Equal(t, int64(5), total1)

	result2, total2, err := testRepository.ListSubscriptions(context.Background(), 2, 2, nil, nil, nil, nil, nil, nil, nil, nil)
	assert.NoError(t, err)
	assert.Len(t, result2, 2)
	assert.Equal(t, int64(5), total2)

	assert.NotEqual(t, result1[0].ID, result2[0].ID)
}

func TestListSubscriptions_Sorting(t *testing.T) {
	SetupRepo(t)
	seedTestSubscriptions(t)

	sortBy := "price"
	sortOrder := "desc"

	result, _, err := testRepository.ListSubscriptions(context.Background(), 1, 10, nil, nil, nil, nil, nil, nil, &sortBy, &sortOrder)

	assert.NoError(t, err)
	assert.True(t, len(result) >= 2)

	for i := 1; i < len(result); i++ {
		assert.GreaterOrEqual(t, result[i-1].Price, result[i].Price)
	}
}

func TestListSubscriptions_EmptyResult(t *testing.T) {
	SetupRepo(t)

	nonExistentUser := "non-existent-user"

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, &nonExistentUser, nil, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Empty(t, result)
	assert.Equal(t, int64(0), total)
}

func TestCalculateTotalCost_Success(t *testing.T) {
	SetupRepo(t)

	sub1 := createTestSubscription("Netflix", "user-1", "01-2025", "12-2025", 1500)
	sub2 := createTestSubscription("Spotify", "user-1", "02-2025", "11-2025", 1200)

	err := testRepository.CreateSubscription(context.Background(), sub1)
	require.NoError(t, err)
	err = testRepository.CreateSubscription(context.Background(), sub2)
	require.NoError(t, err)

	userID := "user-1"
	startDate, _ := time.Parse("01-2006", "01-2025")
	endDate, _ := time.Parse("01-2006", "12-2025")

	result, err := testRepository.CalculateTotalCost(context.Background(), userID, "", &startDate, &endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2700), result)
}

func TestCalculateTotalCost_FilterByService(t *testing.T) {
	SetupRepo(t)

	sub1 := createTestSubscription("Netflix", "user-1", "01-2025", "12-2025", 1500)
	sub2 := createTestSubscription("Spotify", "user-1", "01-2025", "12-2025", 1200)

	err := testRepository.CreateSubscription(context.Background(), sub1)
	require.NoError(t, err)
	err = testRepository.CreateSubscription(context.Background(), sub2)
	require.NoError(t, err)

	userID := "user-1"
	serviceName := "Netflix"
	startDate, _ := time.Parse("01-2006", "01-2025")
	endDate, _ := time.Parse("01-2006", "12-2025")

	result, err := testRepository.CalculateTotalCost(context.Background(), userID, serviceName, &startDate, &endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1500), result)
}

func TestCalculateTotalCost_NoMatches(t *testing.T) {
	SetupRepo(t)

	nonExistentUser := "non-existent-user"
	startDate := time.Now()
	endDate := time.Now().AddDate(0, 1, 0)

	result, err := testRepository.CalculateTotalCost(context.Background(), nonExistentUser, "Netflix", &startDate, &endDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(0), result)
}

func TestRepository_ConcurrentAccess(t *testing.T) {
	SetupRepo(t)

	done := make(chan bool, 2)

	go func() {
		sub := createTestSubscription("Service1", "user-1", "01-2025", "12-2025", 1000)
		err := testRepository.CreateSubscription(context.Background(), sub)
		assert.NoError(t, err)
		done <- true
	}()

	go func() {
		sub := createTestSubscription("Service2", "user-2", "01-2025", "12-2025", 2000)
		err := testRepository.CreateSubscription(context.Background(), sub)
		assert.NoError(t, err)
		done <- true
	}()

	<-done
	<-done

	result, total, err := testRepository.ListSubscriptions(context.Background(), 1, 10, nil, nil, nil, nil, nil, nil, nil, nil)

	assert.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, result, 2)
}

func TestRepository_CancelledContext(t *testing.T) {
	SetupRepo(t)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	sub := createTestSubscription("Netflix", "user-123", "07-2025", "12-2025", 1500)
	err := testRepository.CreateSubscription(ctx, sub)

	assert.Error(t, err)
}
