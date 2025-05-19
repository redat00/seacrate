package system

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/api/middlewares"
	"github.com/redat00/seacrate/internal/models"
	"github.com/redat00/seacrate/internal/shamir"
)

var submittedKeyParts []string

type bodySubmitUnsealPart struct {
	Part string `json:"part" validate:"required"`
}

type responseSubmitUnsealPart struct {
	Result string `json:"result"`
}

func SubmitUnsealPart(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*middlewares.Engines)

	body := new(bodySubmitUnsealPart)
	errs := ctx.Bind().JSON(body)
	if errs != nil {
		ctx.Status(500).JSON(models.ErrorResponse{
			Error:   "Internal Server Error",
			Details: fmt.Sprintf("Error during validation : %v", errs.Error()),
		})
		return nil
	}

	if !engines.EncryptionEngine.GetSealStatus() {
		ctx.Status(400).JSON(models.ErrorResponse{
			Error:   "System is already unsealed.",
			Details: "Seacrate is currently unsealed, but you're trying to unseal it, this is not permitted.",
		})
		return nil
	}

	// Get configuration threshold
	thresholdCountMeta, err := engines.DatabaseEngine.GetMeta("thresholdCount")
	if err != nil {
		ctx.Status(500).JSON(models.ErrorResponse{
			Error:   "Internal Server Error",
			Details: fmt.Sprintf("Could not get meta `thresholdCount` from database : %s", err.Error()),
		})
		return nil
	}

	thresholdCount, err := strconv.Atoi(thresholdCountMeta.Value)
	if err != nil {
		ctx.Status(500).JSON(models.ErrorResponse{
			Error:   "Internal Server Error",
			Details: "Could not convert `thresholdCount` from string to int. Database data might be wrong.",
		})
		return nil
	}

	submittedKeyParts = append(submittedKeyParts, body.Part)
	if len(submittedKeyParts) >= thresholdCount {
		var convertedParts [][]byte
		for i := range submittedKeyParts {
			decodedString, _ := base64.StdEncoding.DecodeString(submittedKeyParts[i])
			convertedParts = append(convertedParts, decodedString)
		}

		key, err := shamir.Combine(convertedParts)
		if err != nil {
			ctx.Status(500).JSON(models.ErrorResponse{
				Error:   "Internal Server Error",
				Details: fmt.Sprintf("Could not combine keys using the Shamir algorithm due to the following error : %s", err.Error()),
			})
			return nil
		}

		engines.EncryptionEngine.SetKey(key)
		engines.EncryptionEngine.SetSealStatus(false)

		submittedKeyParts = make([]string, 0)

		ctx.JSON(models.MessageResponse{
			Details: "Enough parts have been provided to trigger the threshold. Seacrate is now unsealed !",
		})
		return nil
	}

	ctx.JSON(models.MessageResponse{
		Details: "Your key has been added to queue for unseal process. Threshold is still not hit.",
	})
	return nil
}

type responseGetSealStatus struct {
	Status bool `json:"status"`
}

func GetSealStatus(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*middlewares.Engines)
	ctx.JSON(responseGetSealStatus{
		Status: engines.EncryptionEngine.GetSealStatus(),
	})
	return nil
}
