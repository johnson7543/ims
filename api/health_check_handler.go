package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnson7543/ims/db"
)

func NewHealthCheckHandler(store *db.Store) *HealthCheckHandler {
	return &HealthCheckHandler{
		store: store,
	}
}

type HealthCheckHandler struct {
	store *db.Store
}

func (h *HealthCheckHandler) HandleHealthCheck(c *fiber.Ctx) error {
	healthCheck, err := h.store.HealthCheck.CheckHealth(c.Context())
	if err != nil {
		return c.JSON(fiber.Map{
			"server_status": "ok",
			"db_status":     err,
			"db_version":    "unknown",
			"status":        "error occurred in db connection",
		})
	}

	return c.JSON(fiber.Map{
		"server_status": "ok",
		"db_status":     "ok",
		"db_version":    healthCheck,
		"status":        "ok",
	})
}
