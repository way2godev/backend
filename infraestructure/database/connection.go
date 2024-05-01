package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
    dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatalf("DATABASE_DSN environment variable not set")
	}

	var l logger.LogLevel
	if os.Getenv("APP_ENV") == "production" {
		l = logger.Error
	} else {
		l = logger.Info
	}
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(l),
	})
    if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
    }

	raw, _ := db.DB()
	raw.Ping()

	raw.SetMaxOpenConns(10)
	raw.SetMaxIdleConns(10)

	log.Print("Database connection established")
    DB = db
}


func GetDB() *gorm.DB {
	return DB
}