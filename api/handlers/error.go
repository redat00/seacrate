package handlers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/internal/models"
)

func ErrorHandler(ctx fiber.Ctx, err error) error {
	ctx.Status(500).JSON(models.ErrorResponse{
		Error:   "Internal Server Error",
		Details: err.Error(),
	})
	return nil
}
