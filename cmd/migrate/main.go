package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"avito_project/internal/app/ds"
	"avito_project/internal/app/dsn"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(
		&ds.Employee{},
		&ds.Organization{},
		&ds.Tender{},
		&ds.Bid{},
		&ds.OrganizationResponsible{},
	)
	if err != nil {
		panic("cant migrate db")
	}
}
