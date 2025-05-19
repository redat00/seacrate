package api

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/redat00/seacrate/api/middlewares"
	"github.com/redat00/seacrate/api/routers/secrets"
	"github.com/redat00/seacrate/api/routers/system"
	"github.com/redat00/seacrate/internal/database"
	"github.com/redat00/seacrate/internal/encryption"
)

const (
	RootPath = "/api/v1"
)

func NewApi(
	encryptionEngine encryption.EncryptionEngine,
	databaseEngine database.DatabaseEngine,
) *fiber.App {
	app := fiber.New()

	app.Use(
		middlewares.EnginesMiddleware(
			encryptionEngine,
			databaseEngine,
		),
	)
	app.Use(recover.New())

	v1 := app.Group(RootPath)

	systemGroup := v1.Group("/system")
	systemGroup.Post("/seal", system.SubmitUnsealPart)
	systemGroup.Get("/seal", system.GetSealStatus)

	secretGroup := v1.Group("/secrets")
	secretGroup.Use(middlewares.SealMiddleware)
	secretGroup.Get("/*", secrets.GetSecret)
	secretGroup.Post("/*", secrets.CreateSecret)
	secretGroup.Delete("/*", secrets.DeleteSecret)

	return app
}
