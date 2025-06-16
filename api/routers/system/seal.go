package system

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/api/middlewares"
	"github.com/redat00/seacrate/internal/helpers"
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

func emptySubmittedKeyParts() {
	submittedKeyParts = make([]string, 0)
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
		return fmt.Errorf("could not get meta `thresholdCount` from database.")
	}

	thresholdCount, err := strconv.Atoi(thresholdCountMeta.Value)
	if err != nil {
		return fmt.Errorf("could not get convert `thresholdCount` from string to int")
	}

	submittedKeyParts = append(submittedKeyParts, body.Part)
	if len(submittedKeyParts) >= thresholdCount {
		defer emptySubmittedKeyParts()
		var convertedParts [][]byte
		for i := range submittedKeyParts {
			decodedString, _ := base64.StdEncoding.DecodeString(submittedKeyParts[i])
			convertedParts = append(convertedParts, decodedString)
		}

		key, err := shamir.Combine(convertedParts)
		if err != nil {
			return fmt.Errorf("Could not combine keys using the Shamir algorithm : %s", err)
		}

		decryptionKeyHash, err := engines.DatabaseEngine.GetMeta("decryptionKeyHash")
		if err != nil {
			return fmt.Errorf("Could not get meta `decryptionKeyHash` from database.")
		}

		splittedDecryptionKeyHash := strings.Split(decryptionKeyHash.Value, "$")
		salt, err := hex.DecodeString(splittedDecryptionKeyHash[0])
		if err != nil {
			return fmt.Errorf("Could not decode salt : %s", err)
		}

		hash, err := hex.DecodeString(splittedDecryptionKeyHash[1])
		if err != nil {
			return fmt.Errorf("Could not decode hash : %s", err)
		}

		result, err := helpers.Compare(key, hash, salt)
		if err != nil {
			return fmt.Errorf("Could not compare hash : %s", err)
		}

		if !result {
			ctx.Status(400).JSON(models.ErrorResponse{
				Error:   "Wrong Key",
				Details: "The key created from shards is not valid.",
			})
			return nil
		}

		// Get master key
		mKey, err := engines.DatabaseEngine.GetMeta("masterKey")
		if err != nil {
			return fmt.Errorf("Could not get meta `masterKey` from database.")
		}

		decodedMasterKey, err := hex.DecodeString(mKey.Value)
		if err != nil {
			return fmt.Errorf("Could not decode master key.")
		}

		engines.EncryptionEngine.SetKey(key)
		masterKey, err := engines.EncryptionEngine.DecryptData(decodedMasterKey)
		if err != nil {
			fmt.Errorf("Could not decrypt master key.")
		}

		engines.EncryptionEngine.SetKey(masterKey)
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
