package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/redat00/seacrate/internal/config"
	"github.com/redat00/seacrate/internal/database"
	"github.com/redat00/seacrate/internal/encryption"
	"github.com/redat00/seacrate/internal/shamir"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initCmd)
}

type genereatedKeys struct {
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

		var keyPartCount int
		fmt.Println("How many key part do you wish to generate ?")
		fmt.Scan(&keyPartCount)

		var thresholdCount int
		fmt.Println("How many key part are required to unseal the instance ?")
		fmt.Scan(&thresholdCount)

		key, err := encryptionEngine.GenerateKey(32)
		if err != nil {
			panic(err)
		}

		splitted, err := shamir.Split(key, keyPartCount, thresholdCount)
		if err != nil {
			panic(err)
		}

		databaseEngine.CreateMeta("initialized", "yes")
		databaseEngine.CreateMeta("thresholdCount", strconv.Itoa(thresholdCount))
		createdKeys := genereatedKeys{}

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
