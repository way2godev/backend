package services

import (
	"log"
	"time"
	"way2go/domain/entities"
	"way2go/infraestructure/database"
)

type MigrateModelsService struct {}

var MigrateModelsServiceInstance MigrateModelsService = MigrateModelsService{}

func (s *MigrateModelsService) MigrateModels() {
	log.Println("Migrating DB models")
	start := time.Now()

	database.DB.AutoMigrate(&entities.Agency{})
	database.DB.AutoMigrate(&entities.Line{})
	database.DB.AutoMigrate(&entities.LineVariant{})
	database.DB.AutoMigrate(&entities.Schedule{})
	database.DB.AutoMigrate(&entities.ScheduleStop{})
	database.DB.AutoMigrate(&entities.Stop{})
	database.DB.AutoMigrate(&entities.Datasource{})

	log.Printf("Migration completed in %v\n", time.Since(start))
}