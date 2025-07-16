package handlers

import (
	"log/slog"

	"github.com/rasadov/subscription-manager/internal/service"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *slog.Logger
}

func NewSubscriptionHandler(service service.SubscriptionService, logger *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{service: service, logger: logger}
}
