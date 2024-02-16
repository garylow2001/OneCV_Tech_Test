package main

import (
	"github.com/garylow2001/OneCV_Tech_Test/api"
	"github.com/garylow2001/OneCV_Tech_Test/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func init() {
	database.LoadEnvVariables()
	database.ConnectToDB()
	database.SyncDB()
}

func main() {
	db := database.GetDB() // Assuming GetDB returns *gorm.DB
	router := gin.Default()

	setUpRouters(router, db)

	router.Run()
}

func setUpRouters(router *gin.Engine, db *gorm.DB) {
	router.POST("/api/register", func(c *gin.Context) {
		api.RegisterHandler(db, c)
	})
	router.POST("/ping", func(c *gin.Context) {
		api.PingHandler(db, c)
	})
}
