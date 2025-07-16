package tests

import (
	"context"
	"testing"

	"github.com/rasadov/subscription-manager/internal/dto"
)

var (
	testServiceName = "test"
	testPrice       = int64(10)
	testUserID      = "test"
	testStartDate   = "2025-01-01"
	testEndDate     = "2025-01-02"
)

// Unit tests

func TestCreateSubscriptionService(t *testing.T) {
	SetupServiceUnitTests(t)
	req := dto.CreateSubscriptionRequest{
		ServiceName: testServiceName,
		Price:       testPrice,
		UserID:      testUserID,
		StartDate:   testStartDate,
		EndDate:     testEndDate,
	}

	_, err := testService.CreateSubscription(context.Background(), req)
	if err != nil {
		t.Errorf("CreateSubscription error: %v", err)
	}
}

func TestUpdateSubscriptionService(t *testing.T) {
	SetupServiceUnitTests(t)
	updatedServiceName, updatedPrice, updatedStartDate, updatedEndDate := "updated", int64(20), "2025-01-03", "2025-01-04"

	req := dto.UpdateSubscriptionRequest{
		ServiceName: &updatedServiceName,
		Price:       &updatedPrice,
		StartDate:   &updatedStartDate,
		EndDate:     &updatedEndDate,
	}

	_, err := testService.UpdateSubscription(context.Background(), 1, req)
	if err != nil {
		t.Errorf("UpdateSubscription error: %v", err)
	}

	updatedSubscription, err := testService.GetSubscription(context.Background(), 1)
	if err != nil {
		t.Errorf("GetSubscription error: %v", err)
	}

	if updatedSubscription.ServiceName != updatedServiceName {
		t.Errorf("ServiceName expected: %s, actual: %s", updatedServiceName, updatedSubscription.ServiceName)
	}
	if updatedSubscription.Price != updatedPrice {
		t.Errorf("Price expected: %d, actual: %d", updatedPrice, updatedSubscription.Price)
	}
	if updatedSubscription.StartDate != updatedStartDate {
		t.Errorf("StartDate expected: %s, actual: %s", updatedStartDate, updatedSubscription.StartDate)
	}
	if updatedSubscription.EndDate != updatedEndDate {
		t.Errorf("EndDate expected: %s, actual: %s", updatedEndDate, updatedSubscription.EndDate)
	}
}

func TestDeleteSubscriptionService(t *testing.T) {
	SetupServiceUnitTests(t)
	if err := testService.DeleteSubscription(context.Background(), 1); err != nil {
		t.Errorf("DeleteSubscription error: %v", err)
	}
}

func TestListSubscriptionsService(t *testing.T) {
	SetupServiceUnitTests(t)
	req := dto.ListSubscriptionsQuery{
		UserID: &testUserID,
		Page:   1,
		Limit:  10,
	}

	_, err := testService.ListSubscriptions(context.Background(), req)
	if err != nil {
		t.Errorf("ListSubscriptions error: %v", err)
	}
}

func TestCalculateTotalCostService(t *testing.T) {
	SetupServiceUnitTests(t)
	req := dto.TotalCostQuery{
		UserID:      &testUserID,
		ServiceName: &testServiceName,
		StartDate:   &testStartDate,
		EndDate:     &testEndDate,
	}

	_, err := testService.CalculateTotalCost(context.Background(), req)
	if err != nil {
		t.Errorf("CalculateTotalCost error: %v", err)
	}
}

// Integration tests

func TestCreateSubscriptionIntegration(t *testing.T) {
	SetupServiceIntegrationTests(t)
	req := dto.CreateSubscriptionRequest{
		ServiceName: testServiceName,
		Price:       testPrice,
		UserID:      testUserID,
		StartDate:   testStartDate,
		EndDate:     testEndDate,
	}

	_, err := testService.CreateSubscription(context.Background(), req)
	if err != nil {
		t.Errorf("CreateSubscription error: %v", err)
	}
}

func TestUpdateSubscriptionIntegration(t *testing.T) {
	SetupServiceIntegrationTests(t)
	updatedServiceName, updatedPrice, updatedStartDate, updatedEndDate := "updated", int64(20), "2025-01-03", "2025-01-04"

	req := dto.UpdateSubscriptionRequest{
		ServiceName: &updatedServiceName,
		Price:       &updatedPrice,
		StartDate:   &updatedStartDate,
		EndDate:     &updatedEndDate,
	}

	_, err := testService.UpdateSubscription(context.Background(), 1, req)
	if err != nil {
		t.Errorf("UpdateSubscription error: %v", err)
	}
}

func TestDeleteSubscriptionIntegration(t *testing.T) {
	SetupServiceIntegrationTests(t)
	if err := testService.DeleteSubscription(context.Background(), 1); err != nil {
		t.Errorf("DeleteSubscription error: %v", err)
	}
}

func TestListSubscriptionsIntegration(t *testing.T) {
	SetupServiceIntegrationTests(t)
	req := dto.ListSubscriptionsQuery{
		UserID: &testUserID,
	}
	_, err := testService.ListSubscriptions(context.Background(), req)
	if err != nil {
		t.Errorf("ListSubscriptions error: %v", err)
	}
}

func TestCalculateTotalCostIntegration(t *testing.T) {
	SetupServiceIntegrationTests(t)
	req := dto.TotalCostQuery{
		UserID:      &testUserID,
		ServiceName: &testServiceName,
		StartDate:   &testStartDate,
		EndDate:     &testEndDate,
	}
	_, err := testService.CalculateTotalCost(context.Background(), req)
	if err != nil {
		t.Errorf("CalculateTotalCost error: %v", err)
	}
}
