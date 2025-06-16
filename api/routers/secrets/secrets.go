package secrets // import github.com/redat00/seacrate/api/routers/secrets

import (
	"encoding/hex"
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/api/middlewares"
	seacrateErrors "github.com/redat00/seacrate/internal/errors"
	"github.com/redat00/seacrate/internal/models"
)

type inputCreateSecret struct {
	Value string `json:"value" validate:"required"`
}

func CreateSecret(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*middlewares.Engines)
	secretKey := ctx.Params("*")

	// We need to append a / to our secretKey
	secretKey = "/" + secretKey

	// Validate body
	in := inputCreateSecret{}
	errs := ctx.Bind().JSON(&in)
	if errs != nil {
		ctx.Status(400).JSON(models.ErrorResponse{
			Error:   "Body invalid",
			Details: errs.Error(),
		})
		return nil
	}

	encryptedValue, err := engines.EncryptionEngine.EncryptData([]byte(in.Value))
	if err != nil {
		return fmt.Errorf("An error happened during encryption of the data.")
	}

	// Create secret in database
	err = engines.DatabaseEngine.CreateSecret(secretKey, hex.EncodeToString(encryptedValue))
	if err != nil {
		switch err.(type) {
		case seacrateErrors.ErrSecretDuplicateKey:
			ctx.Status(400).JSON(models.ErrorResponse{
				Error:   "Could not create secret",
				Details: err.Error(),
			})
			return nil
		case seacrateErrors.ErrOverridingFolder:
			ctx.Status(400).JSON(models.ErrorResponse{
				Error:   "Could not create secret",
				Details: err.Error(),
			})
			return nil
		case seacrateErrors.ErrOverridingSecret:
			ctx.Status(400).JSON(models.ErrorResponse{
				Error:   "Could not create secret",
				Details: err.Error(),
			})
			return nil
		default:
			return fmt.Errorf("An error happened during creation of the secret")
		}
	}

	// Return a succesfull response
	ctx.JSON(models.MessageResponse{
		Details: "The secret was successfully created",
	})
	return nil
}

type responseGetSecretSecret struct {
	Type   string         `json:"type"`
	Secret *models.Secret `json:"secret"`
}

type responseGetSecretFolder struct {
	Type    string                 `json:"type"`
	Content []models.FolderContent `json:"content"`
}

func GetSecret(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*middlewares.Engines)
	secretKey := ctx.Params("*")

	// We need to append a / to our secretKey
	secretKey = "/" + secretKey

	multiple, secrets, secret, err := engines.DatabaseEngine.GetSecret(secretKey)
	if err != nil {
		switch err.(type) {
		case seacrateErrors.ErrSecretNotFound:
			ctx.Status(404).JSON(models.ErrorResponse{
				Error:   "Secret not found",
				Details: err.Error(),
			})
			return nil
		default:
			return fmt.Errorf("An error happened while trying to get the secret.")
		}
	}

	if multiple {
		if secrets == nil {
			secrets = []models.FolderContent{}
		}
		ctx.JSON(responseGetSecretFolder{
			Type:    "folder",
			Content: secrets,
		})
		return nil
	}

	stringToHex, err := hex.DecodeString(secret.Value)
	if err != nil {
		return fmt.Errorf("An error happened while decoding hex to bytes.")
	}

	decryptedData, err := engines.EncryptionEngine.DecryptData(stringToHex)
	if err != nil {
		return fmt.Errorf("An error happened during decryption of the data.")
	}

	secret.Value = string(decryptedData)
	ctx.JSON(responseGetSecretSecret{
		Type:   "secret",
		Secret: secret,
	})
	return nil
}

func DeleteSecret(ctx fiber.Ctx) error {
	engines := ctx.Locals("engines").(*middlewares.Engines)
	secretKey := ctx.Params("*")

	// We need to append a / to our secretKey
	secretKey = "/" + secretKey

	err := engines.DatabaseEngine.DeleteSecret(secretKey)
	if err != nil {
		switch err.(type) {
		case seacrateErrors.ErrSecretNotFound:
			ctx.Status(404).JSON(models.ErrorResponse{
				Error:   "Secret not found",
				Details: err.Error(),
			})
			return nil
		default:
			return fmt.Errorf("An error happened during the deletion of the secret.")
		}
	}

	ctx.JSON(models.MessageResponse{
		Details: "The secret was successfully deleted",
	})
	return nil
}
