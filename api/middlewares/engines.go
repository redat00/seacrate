package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/internal/database"
	"github.com/redat00/seacrate/internal/encryption"
)

type Engines struct {
	EncryptionEngine encryption.EncryptionEngine
	DatabaseEngine   database.DatabaseEngine
}

func EnginesMiddleware(
	encrEngine encryption.EncryptionEngine,
	dbEngine database.DatabaseEngine,
) func(ctx fiber.Ctx) error {
	engines := Engines{
		encrEngine,
		dbEngine,
	}
	return func(ctx fiber.Ctx) error {
		ctx.Locals("engines", &engines)
		return ctx.Next()
	}
}
