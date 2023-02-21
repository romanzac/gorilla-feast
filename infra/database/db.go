package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB Instance
var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	// Connect to Postgres
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

}
