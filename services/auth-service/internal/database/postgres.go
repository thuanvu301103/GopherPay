package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes a connection to any supported RDBMS
func InitPostgres(dbURL string) (*gorm.DB, error) {
	var dialector gorm.Dialector

	dialector = postgres.Open(dbURL)

	return gorm.Open(dialector, &gorm.Config{})
}
