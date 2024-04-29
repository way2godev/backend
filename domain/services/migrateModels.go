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
	database.DB.AutoMigrate(&entities.Calendar{})
	database.DB.AutoMigrate(&entities.CalendarException{})
	database.DB.AutoMigrate(&entities.Datasource{})
	database.DB.AutoMigrate(&entities.Line{})
	database.DB.AutoMigrate(&entities.Schedule{})
	database.DB.AutoMigrate(&entities.ScheduleStop{})
	database.DB.AutoMigrate(&entities.Shape{})
	database.DB.AutoMigrate(&entities.ShapeElement{})
	database.DB.AutoMigrate(&entities.Stop{})

	log.Printf("Migration completed in %v\n", time.Since(start))
}