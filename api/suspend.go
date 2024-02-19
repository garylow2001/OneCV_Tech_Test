package api

import (
	"fmt"
	"net/http"

	"github.com/garylow2001/OneCV_Tech_Test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SuspendRequest struct {
	StudentEmail string `json:"student" binding:"required"`
}

func SuspendHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestData SuspendRequest
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Get student by email
		studentEmail := requestData.StudentEmail
		student, getStudentErr := getStudent(db, studentEmail)
		if getStudentErr != nil {
			c.JSON(http.StatusNotFound, gin.H{"message": fmt.Sprintf("Student %s not found", studentEmail)})
			return
		}

		// Check if student is already suspended
		if student.Suspended {
			c.JSON(http.StatusConflict, gin.H{"message": fmt.Sprintf("Student %s is already suspended", studentEmail)})
			return
		}

		// Suspend student
		db.Model(&student).Update("Suspended", true)

		c.Status(http.StatusNoContent)
	}
}

func getStudent(db *gorm.DB, studentEmail string) (*models.Student, error) {
	var student models.Student
	result := db.Where("email = ?", studentEmail).First(&student)
	if result.Error != nil {
		return nil, result.Error
	}
	return &student, nil
}
