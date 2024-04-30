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

func (s *MigrateModelsService) DropModels() {
	log.Println("Dropping DB models")
	start := time.Now()

	database.DB.Migrator().DropTable(&entities.Agency{})
	database.DB.Migrator().DropTable(&entities.Calendar{})
	database.DB.Migrator().DropTable(&entities.CalendarException{})
	database.DB.Migrator().DropTable(&entities.Datasource{})
	database.DB.Migrator().DropTable(&entities.Line{})
	database.DB.Migrator().DropTable(&entities.Schedule{})
	database.DB.Migrator().DropTable(&entities.ScheduleStop{})
	database.DB.Migrator().DropTable(&entities.Shape{})
	database.DB.Migrator().DropTable(&entities.ShapeElement{})
	database.DB.Migrator().DropTable(&entities.Stop{})

	log.Printf("Drop completed in %v\n", time.Since(start))
}