package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/pkg/exceptions"
	"github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/service"
	"github.com/rasadov/subscription-manager/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestCreateSubscriptionService_Success(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       1500,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "07-2025",
		EndDate:     "12-2025",
	}

	mockRepo.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Run(func(args mock.Arguments) {
			sub := args.Get(1).(*models.Subscription)
			sub.ID = 1
		}).Return(nil)

	result, err := service.CreateSubscription(context.Background(), req)

	startDateParsed, _ := time.Parse("01-2006", req.StartDate)
	endDateParsed, _ := time.Parse("01-2006", req.EndDate)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, req.ServiceName, result.ServiceName)
	assert.Equal(t, req.Price, result.Price)
	assert.Equal(t, req.UserID, result.UserID)
	assert.Equal(t, startDateParsed, time.Time(result.StartDate))
	if assert.NotNil(t, result.EndDate) {
		assert.Equal(t, endDateParsed, time.Time(*result.EndDate))
	}
	assert.Equal(t, uint(1), result.ID)

	mockRepo.AssertExpectations(t)
}

func TestCreateSubscriptionService_RepositoryError(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Netflix",
		Price:       1500,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "07-2025",
	}

	expectedError := errors.New("database connection failed")
	mockRepo.On("CreateSubscription", mock.Anything, mock.AnythingOfType("*models.Subscription")).
		Return(expectedError)

	result, err := service.CreateSubscription(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Implements(t, (*exceptions.HTTPError)(nil), err)
	assert.Contains(t, err.Error(), expectedError.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetSubscriptionService_Success(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	startDateParsed, _ := time.Parse("01-2006", "07-2025")
	endDateParsed, _ := time.Parse("01-2006", "12-2025")

	expectedSub := &models.Subscription{
		ID:          1,
		ServiceName: "Netflix",
		Price:       1500,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   startDateParsed,
		EndDate:     &endDateParsed,
	}

	mockRepo.On("GetSubscription", mock.Anything, 1).Return(expectedSub, nil)

	result, err := service.GetSubscription(context.Background(), 1)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedSub.ID, result.ID)
	assert.Equal(t, expectedSub.ServiceName, result.ServiceName)
	assert.Equal(t, expectedSub.Price, result.Price)
	assert.Equal(t, expectedSub.StartDate, time.Time(result.StartDate))
	if assert.NotNil(t, result.EndDate) && expectedSub.EndDate != nil {
		assert.Equal(t, *expectedSub.EndDate, time.Time(*result.EndDate))
	}

	mockRepo.AssertExpectations(t)
}

func TestGetSubscriptionService_NotFound(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	mockRepo.On("GetSubscription", mock.Anything, 999).Return(nil, gorm.ErrRecordNotFound)

	result, err := service.GetSubscription(context.Background(), 999)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Implements(t, (*exceptions.HTTPError)(nil), err)
	assert.Contains(t, err.Error(), gorm.ErrRecordNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestUpdateSubscriptionService_Success(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	startDateParsed, _ := time.Parse("01-2006", "07-2025")

	existingSub := &models.Subscription{
		ID:          1,
		ServiceName: "Netflix",
		Price:       1500,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   startDateParsed,
	}

	newServiceName := "Spotify"
	newPrice := int64(1200)

	req := dto.UpdateSubscriptionRequest{
		ServiceName: &newServiceName,
		Price:       &newPrice,
	}

	mockRepo.On("GetSubscription", mock.Anything, 1).Return(existingSub, nil)
	mockRepo.On("UpdateSubscription", mock.Anything, 1, mock.AnythingOfType("*models.Subscription")).Return(nil)

	result, err := service.UpdateSubscription(context.Background(), 1, req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, newServiceName, result.ServiceName)
	assert.Equal(t, newPrice, result.Price)
	assert.Equal(t, existingSub.UserID, result.UserID)

	mockRepo.AssertExpectations(t)
}

func TestListSubscriptionsService_Success(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	expectedSubs := []*models.Subscription{
		{ID: 1, ServiceName: "Netflix", Price: 1500},
		{ID: 2, ServiceName: "Spotify", Price: 1200},
	}

	query := dto.ListSubscriptionsQuery{
		Page:  1,
		Limit: 10,
	}

	mockRepo.On(
		"ListSubscriptions",
		mock.Anything,           // ctx
		query.Page,                // page
		query.Limit,              // elements
		(*string)(nil),          // userID
		(*string)(nil),          // serviceName
		(*time.Time)(nil),       // startDateFrom
		(*time.Time)(nil),       // startDateTo
		(*time.Time)(nil),       // endDateFrom
		(*time.Time)(nil),       // endDateTo
		mock.AnythingOfType("*string"), // sortBy
		mock.AnythingOfType("*string"), // sortOrder
	).Return(expectedSubs, int64(2), nil)

	result, err := service.ListSubscriptions(context.Background(), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Data, 2)
	assert.Equal(t, 2, result.Pagination.Total)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.Limit)
	assert.Equal(t, 1, result.Pagination.TotalPages)

	mockRepo.AssertExpectations(t)
}

func TestListSubscriptionsService_DefaultPagination(t *testing.T) {
	mockRepo := new(mocks.SubscriptionRepository)
	service := service.NewSubscriptionService(mockRepo)

	query := dto.ListSubscriptionsQuery{}

	expectedQuery := query
	expectedQuery.Page = 1
	expectedQuery.Limit = 10

	mockRepo.On(
		"ListSubscriptions",
		mock.Anything,           // ctx
		expectedQuery.Page,         // page
		expectedQuery.Limit,        // elements
		(*string)(nil),          // userID
		(*string)(nil),          // serviceName
		(*time.Time)(nil),       // startDateFrom
		(*time.Time)(nil),       // startDateTo
		(*time.Time)(nil),       // endDateFrom
		(*time.Time)(nil),       // endDateTo
		mock.AnythingOfType("*string"), // sortBy
		mock.AnythingOfType("*string"), // sortOrder
	).Return([]*models.Subscription{}, int64(0), nil)

	result, err := service.ListSubscriptions(context.Background(), query)

	assert.NoError(t, err)
	assert.Equal(t, 1, result.Pagination.Page)
	assert.Equal(t, 10, result.Pagination.Limit)

	mockRepo.AssertExpectations(t)
}

func TestCreateAndGetSubscriptionIntegration(t *testing.T) {
	SetupRepo(t)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Integration Test Service",
		Price:       2500,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "07-2025",
		EndDate:     "12-2025",
	}

	created, err := testService.CreateSubscription(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, created)
	assert.NotZero(t, created.ID)

	retrieved, err := testService.GetSubscription(context.Background(), int(created.ID))
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, created.ID, retrieved.ID)
	assert.Equal(t, req.ServiceName, retrieved.ServiceName)
	assert.Equal(t, req.Price, retrieved.Price)
	assert.Equal(t, req.UserID, retrieved.UserID)
}

func TestFullCRUDIntegration(t *testing.T) {
	SetupRepo(t)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "CRUD Test Service",
		Price:       1800,
		UserID:      "123e4567-e89b-12d3-a456-426614174000",
		StartDate:   "07-2025",
	}

	created, err := testService.CreateSubscription(context.Background(), req)
	assert.NoError(t, err)
	assert.NotZero(t, created.ID)

	newServiceName := "Updated Service"
	updateReq := dto.UpdateSubscriptionRequest{
		ServiceName: &newServiceName,
	}

	updated, err := testService.UpdateSubscription(context.Background(), int(created.ID), updateReq)
	assert.NoError(t, err)
	assert.Equal(t, newServiceName, updated.ServiceName)
	assert.Equal(t, req.Price, updated.Price)

	listQuery := dto.ListSubscriptionsQuery{
		UserID: &req.UserID,
	}

	listResult, err := testService.ListSubscriptions(context.Background(), listQuery)
	assert.NoError(t, err)
	assert.True(t, listResult.Pagination.Total >= 1)

	err = testService.DeleteSubscription(context.Background(), int(created.ID))
	assert.NoError(t, err)

	_, err = testService.GetSubscription(context.Background(), int(created.ID))
	assert.Error(t, err)
	assert.Implements(t, (*exceptions.HTTPError)(nil), err)
assert.Contains(t, err.Error(), gorm.ErrRecordNotFound.Error())
}
