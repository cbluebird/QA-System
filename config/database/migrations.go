package database

import (
	"QA-System/app/models"

	"gorm.io/gorm"
)

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Survey{},
		&models.Question{},
		&models.Option{},
		&models.Manage{},
	)
}
