package cmd

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/redat00/seacrate/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(validateCmd)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file is valid",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration from file
		_, err := config.GenerateConfigFromFile("config.yml")
		if err != nil {
			errors := err.(validator.ValidationErrors)
			fmt.Println(errors)
			os.Exit(1)
		}
	},
}
