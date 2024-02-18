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
	db := database.GetDB()
	router := gin.Default()

	setUpRouters(router, db)

	router.Run()
}

func setUpRouters(router *gin.Engine, db *gorm.DB) {
	// Decouple db from the handlers for testability
	router.POST("/api/register", api.RegisterHandler(db))
	router.GET("/api/commonstudents", api.CommonStudentsHandler(db))
	router.POST("/api/suspend", api.SuspendHandler(db))
	router.POST("/api/retrievefornotifications", api.RetrieveForNotificationsHandler(db))
}
