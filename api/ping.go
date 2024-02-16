package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func PingHandler(db *gorm.DB, c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
