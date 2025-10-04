package database

import (
	"recipe/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Recipe{}, &models.Ingredient{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
