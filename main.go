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
	router.POST("/api/register", api.RegisterHandler(db))
	router.GET("/api/commonstudents", api.CommonStudentsHandler(db))
	router.POST("/api/suspend", api.SuspendHandler(db))
	router.POST("/api/retrievefornotifications", api.RetrieveForNotificationsHandler(db))
}
