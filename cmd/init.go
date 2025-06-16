package cmd

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redat00/seacrate/internal/config"
	"github.com/redat00/seacrate/internal/database"
	"github.com/redat00/seacrate/internal/encryption"
	"github.com/redat00/seacrate/internal/helpers"
	"github.com/redat00/seacrate/internal/shamir"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

func initializeEngines(config *config.Config) (database.DatabaseEngine, encryption.EncryptionEngine) {
	// Create a database engine
	databaseEngine, err := database.NewDatabaseEngine(config.Database)
	if err != nil {
		panic(err)
	}

	// Create an encryption engine
	encryptionEngine, err := encryption.NewEncryptionEngine(config.Encryption)
	if err != nil {
		panic(err)
	}

	return databaseEngine, encryptionEngine
}

type generatedKeys struct {
	Keys []string `json:"keys"`
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Init the application database and create an encryption key",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration from file
		config, err := config.GenerateConfigFromFile("config.yml")
		if err != nil {
			panic(err)
		}

		// Create engines
		databaseEngine, encryptionEngine := initializeEngines(config)

		var keyPartCount int
		fmt.Println("How many key part do you wish to generate ?")
		fmt.Scan(&keyPartCount)

		var thresholdCount int
		fmt.Println("How many key part are required to unseal the instance ?")
		fmt.Scan(&thresholdCount)

		// Generate the master key
		masterKey, err := encryptionEngine.GenerateKey(32)
		if err != nil {
			panic(err)
		}

		// Generate the decryption key
		decryptionKey, err := encryptionEngine.GenerateKey(32)
		if err != nil {
			panic(err)
		}

		// Use the decryption key to encrypt the master key
		encryptionEngine.SetKey(decryptionKey)
		encryptedMasterKey, err := encryptionEngine.EncryptData(masterKey)
		if err != nil {
			panic(err)
		}

		// Obtain decryption key hash
		hash, salt, err := helpers.GenerateHash(decryptionKey, []byte{})
		if err != nil {
			panic(err)
		}

		splitted, err := shamir.Split(decryptionKey, keyPartCount, thresholdCount)
		if err != nil {
			panic(err)
		}

		databaseEngine.CreateMeta("initialized", "yes")
		databaseEngine.CreateMeta("thresholdCount", strconv.Itoa(thresholdCount))
		databaseEngine.CreateMeta("decryptionKeyHash", fmt.Sprintf("%s$%s", hex.EncodeToString(salt), hex.EncodeToString(hash)))
		databaseEngine.CreateMeta("masterKey", hex.EncodeToString(encryptedMasterKey))
		createdKeys := generatedKeys{}

		for _, element := range splitted {
			createdKeys.Keys = append(createdKeys.Keys, base64.StdEncoding.EncodeToString(element))
		}

		rawCreatedKeys, err := json.Marshal(createdKeys)
		if err != nil {
			panic(err)
		}

		err = os.WriteFile("results.json", rawCreatedKeys, 0644)
		if err != nil {
			panic(err)
		}

		fmt.Println("Create file `results.json`")
	},
}
