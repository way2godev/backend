package bootstrap

import (
	"os"
	"way2go/domain/services"
	"way2go/infraestructure/cache"
	"way2go/infraestructure/database"

	"github.com/joho/godotenv"
)

func Init() {
	// This is a placeholder for the bootstrap process
	godotenv.Load()

	// Initialize the database
	database.InitDB()
	cache.InitCache()

	// Migrate the models if the --migrate flag is set on argv or the app is in production
	if len(os.Args) > 1 && os.Args[1] == "--migrate" || os.Getenv("APP_ENV") == "production" {
		services.MigrateModelsServiceInstance.MigrateModels()
	} else if len(os.Args) > 1 && os.Args[1] == "--drop" {
		services.MigrateModelsServiceInstance.DropModels()
		services.MigrateModelsServiceInstance.MigrateModels()
	}
}