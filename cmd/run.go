package cmd

import (
	"github.com/gofiber/fiber/v3"
	"github.com/redat00/seacrate/api"
	"github.com/redat00/seacrate/internal/config"
	"github.com/redat00/seacrate/internal/database"
	"github.com/redat00/seacrate/internal/encryption"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start the server",
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

		seacrateApi := api.NewApi(encryptionEngine, databaseEngine)
		seacrateApi.Listen(":3000", fiber.ListenConfig{
			DisableStartupMessage: true,
		})
	},
}
