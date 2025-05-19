package middlewares

import "github.com/gofiber/fiber/v3"

type responseSealMiddleware struct {
	Error string `json:"error"`
}

func SealMiddleware(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*Engines)
	if status := engines.EncryptionEngine.GetSealStatus(); status {
		resp := responseSealMiddleware{
			Error: "Seacrate is currently sealed. This operation is not permitted.",
		}
		ctx.Status(403)
		return ctx.JSON(resp)
	}
	return ctx.Next()
}
