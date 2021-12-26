package app

import (
	"log"
	"os"

	"go-get-it/config"
	"go-get-it/seeds"

	"github.com/joho/godotenv"
)

func Run() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	logger := &Log{}
	logger.InitialLog()

	InfoLog.Println("Starting the application...")

	config := config.GetConfig()
	app := &App{}
	app.Initialize(config, InfoLog)

	// seeding data to database
	if config.Seed.ClearDB || config.Seed.PopulateDb {
		seeds.Run(app.DB, *config.Seed)
	} else {
		app.Run(":" + os.Getenv("PORT"))
	}

}
