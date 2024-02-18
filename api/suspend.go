package api

import (
	"fmt"
	"net/http"

	"github.com/garylow2001/OneCV_Tech_Test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SuspendHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestData struct {
			StudentEmail string `json:"student" binding:"required"`
		}
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Check if student exists
		var student models.Student
		result := db.Where("email = ?", requestData.StudentEmail).First(&student)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Student %s not found", requestData.StudentEmail)})
			return
		}

		// Check if student is already suspended
		if student.Suspended {
			c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("Student %s is already suspended", requestData.StudentEmail)})
			return
		}

		// Suspend student
		db.Model(&student).Update("Suspended", true)

		c.Status(http.StatusNoContent)
	}
}
