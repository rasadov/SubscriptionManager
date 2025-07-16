package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/subscription-manager/internal/dto"
	_ "github.com/rasadov/subscription-manager/internal/models"
	"github.com/rasadov/subscription-manager/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *slog.Logger
}

func NewSubscriptionHandler(service service.SubscriptionService, logger *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

// CreateSubscription godoc
// @Summary Create a new subscription
// @Description Create a new subscription with the provided details
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body dto.CreateSubscriptionRequest true "Subscription details"
// @Success 201 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [post]
func (h *SubscriptionHandler) CreateSubscription(c *gin.Context) {
	var req dto.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CreateSubscription(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create subscription", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create subscription"})
		return
	}

	h.logger.Info("Subscription created successfully", "id", response.ID)
	c.JSON(http.StatusCreated, response)
}

// GetSubscription godoc
// @Summary Get a subscription by ID
// @Description Get subscription details by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetSubscription(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("Invalid subscription ID", "id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	response, err := h.service.GetSubscription(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to get subscription", "id", id, "error", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Subscription not found"})
		return
	}

	h.logger.Info("Subscription retrieved successfully", "id", id)
	c.JSON(http.StatusOK, response)
}

// UpdateSubscription godoc
// @Summary Update a subscription
// @Description Update subscription details by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body dto.UpdateSubscriptionRequest true "Updated subscription details"
// @Success 200 {object} models.Subscription
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) UpdateSubscription(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("Invalid subscription ID", "id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.UpdateSubscription(c.Request.Context(), id, req)
	if err != nil {
		h.logger.Error("Failed to update subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update subscription"})
		return
	}

	h.logger.Info("Subscription updated successfully", "id", id)
	c.JSON(http.StatusOK, response)
}

// DeleteSubscription godoc
// @Summary Delete a subscription
// @Description Delete a subscription by ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) DeleteSubscription(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		h.logger.Error("Invalid subscription ID", "id", idParam)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid subscription ID"})
		return
	}

	err = h.service.DeleteSubscription(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("Failed to delete subscription", "id", id, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete subscription"})
		return
	}

	h.logger.Info("Subscription deleted successfully", "id", id)
	c.JSON(http.StatusNoContent, nil)
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description Get a list of subscriptions with optional filtering and pagination
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param start_date_from query string false "Start date from filter (MM-YYYY)"
// @Param end_date_from query string false "End date from filter (MM-YYYY)"
// @Param end_date_to query string false "End date to filter (MM-YYYY)"
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order (asc/desc)"
// @Success 200 {object} dto.ListSubscriptionsResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions [get]
func (h *SubscriptionHandler) ListSubscriptions(c *gin.Context) {
	var query dto.ListSubscriptionsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error("Invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.ListSubscriptions(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to list subscriptions", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list subscriptions"})
		return
	}

	h.logger.Info("Subscriptions listed successfully", "count", response.Pagination.Total)
	c.JSON(http.StatusOK, response)
}

// CalculateTotalCost godoc
// @Summary Calculate total cost
// @Description Calculate total cost of subscriptions for a given period with optional filters
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "User ID filter"
// @Param service_name query string false "Service name filter"
// @Param start_date query string true "Start date (MM-YYYY)"
// @Param end_date query string true "End date (MM-YYYY)"
// @Success 200 {object} dto.TotalCostResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /subscriptions/total-cost [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	var query dto.TotalCostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		h.logger.Error("Invalid query parameters", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.CalculateTotalCost(c.Request.Context(), query)
	if err != nil {
		h.logger.Error("Failed to calculate total cost", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate total cost"})
		return
	}

	h.logger.Info("Total cost calculated successfully", "total_cost", response.TotalCost)
	c.JSON(http.StatusOK, response)
}
