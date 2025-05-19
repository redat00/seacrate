package main

import (
	"github.com/redat00/seacrate/cmd"
)

func main() {
	//	encryptionConfig := config.EncryptionConfiguration{
	//		EncryptionAlgorithm: "aes",
	//	}
	//	encryptionEngine, err := encryption.NewEncryptionEngine(encryptionConfig)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	databaseEngine := database.NewDatabaseEngine()
	//	app := api.NewApi(
	//		encryptionEngine,
	//		databaseEngine,
	//	)
	//	app.Listen(":3000")
	//config, err := config.GenerateConfigFromFile("config.yml")
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%+v", config)
	cmd.Execute()
}
