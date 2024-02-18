package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := "host=rain.db.elephantsql.com user=rjyryzuq password=DNMNBz_XZT5m2C-EoeMTyBpVos0n-doa dbname=rjyryzuq port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if DB == nil {
		panic("database connection is nil")
	}

	if err != nil {
		panic("failed to connect database")
	}
}
