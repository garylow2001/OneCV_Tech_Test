package database

import (
	"github.com/garylow2001/OneCV_Tech_Test/models"
	"gorm.io/gorm"
)

func SyncDB() {
	// Auto migrate all models
	err := DB.AutoMigrate(&models.Teacher{}, &models.Student{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return DB
}
