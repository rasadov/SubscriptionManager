package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rasadov/subscription-manager/internal/dto"
	"github.com/rasadov/subscription-manager/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *slog.Logger
}

func NewSubscriptionHandler(service service.SubscriptionService, logger *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}

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

	c.JSON(http.StatusCreated, response)
}

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

	c.JSON(http.StatusOK, response)
}

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

	c.JSON(http.StatusOK, response)
}

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

	c.JSON(http.StatusNoContent, nil)
}

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

	c.JSON(http.StatusOK, response)
}

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

	c.JSON(http.StatusOK, response)
}
