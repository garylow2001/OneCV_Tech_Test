package api

import (
	"net/http"

	"github.com/garylow2001/OneCV_Tech_Test/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Parse request body
		var requestData struct {
			Teacher  string   `json:"teacher" binding:"required"`
			Students []string `json:"students" binding:"required"`
		}
		if err := c.BindJSON(&requestData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		// Check if teacher exists
		var teacher models.Teacher
		result := db.Where("email = ?", requestData.Teacher).First(&teacher)
		if result.Error != nil {
			// Create teacher if not found
			teacher = models.Teacher{Email: requestData.Teacher}
			createTeacherError := db.Create(&teacher)
			if createTeacherError.Error != nil {
				c.JSON(http.StatusBadRequest, gin.H{"message": createTeacherError.Error})
				return
			}
		}

		// Register students
		for _, studentEmail := range requestData.Students {
			var student models.Student
			result := db.Where("email = ?", studentEmail).First(&student)
			if result.Error != nil {
				// Create student if not found
				student = models.Student{Email: studentEmail}
				createStudentError := db.Create(&student)
				if createStudentError.Error != nil {
					c.JSON(http.StatusBadRequest, gin.H{"message": createStudentError.Error})
					return
				}
			}
			// Associate student with teacher
			db.Model(&teacher).Association("Students").Append(&student)
		}

		c.Status(http.StatusNoContent)
	}
}
